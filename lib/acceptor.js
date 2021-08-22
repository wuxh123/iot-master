const fs = require("fs");
const path = require("path");
const EventEmitter = require("events");

const tunnel = require("./tunnel");
const device = require("./device");
const project = require("./project");
const {createEvent} = require("./event");

const mongo = require_plugin("mongodb");

class Acceptor extends EventEmitter {

}

const acceptors = {};

exports.getAcceptor = function (id) {
    return acceptors[id];
}

exports.removeAcceptor = function (id) {
    const acc = acceptors[id];
    if (acc) {
        acc.close();
        delete acc[id];  //projects[id]=null
    }
}

const tunnels = {};

exports.getTunnel = function (id) {
    return tunnels[id];
}

exports.removeTunnel = function (id) {
    const tunnel = tunnels[id];
    if (tunnel) {
        tunnel.close();
        delete tunnel[id];  //projects[id]=null
    }
}


/**
 * 创建接收器（服务）
 * @param model 参数
 * @return {Acceptor}
 */
exports.create = function (model) {
    //检查js脚本是否存在
    const mod = path.join(__dirname, '..', 'acceptors', model.type + '.js');
    if (!fs.existsSync(mod)) {
        throw new Error("不支持的接收器类型：" + model.type);
    }

    const acceptor = require(mod)(model);

    //保存到全局
    acceptors[model._id] = acceptor;

    acceptor.on('connect', async (sn, socket) => {
        //根据SN找通道记录
        let modTunnel = await mongo.db.collection("tunnel").findOne({acceptor_id: model._id, sn},);
        if (modTunnel) {
            if (!modTunnel.enable) {
                socket.write("denied sn")
                return;
            }
            await mongo.db.collection("tunnel").updateOne({_id: modTunnel._id},{$set:{last: new Date(), remote: socket.remoteAddress}});

        } else {
            modTunnel = {acceptor_id: model._id, sn, adapter: model.adapter, enable: true, last: new Date(), remote: socket.remoteAddress};
            const ret = await mongo.db.collection("tunnel").insertOne(modTunnel)
            modTunnel._id = ret.insertedId;

            //TODO 创建设备

        }

        //添加上线记录
        createEvent({tunnel_id: modTunnel._id, event: "上线"})

        //打开通道
        tunnel.open(socket, modTunnel);

        //首次上线
        if (!modTunnel.hasOwnProperty('enable')) {
            //补充启用
            modTunnel.enable = true;
            mongo.db.collection("tunnel").updateOne({_id: modTunnel._id}, {$set: {enable: true}}).then();

            //创建设备，存入数据库
            if (tunnel2.adapter && model.devices && model.devices.length) {
                //依次创建设备
                const modDevices = model.devices.map(d => {
                    return {
                        tunnel_id: modTunnel._id,
                        element_id: d.element_id,
                        slave: d.slave,
                        enable: true,
                    };
                })
                await mongo.db.collection("device").insertMany(modDevices);
            }
        }


        //保存全局通道
        tunnels[modTunnel._id] = tunnel2

        //TODO 此处解析通道配置，初始化通道协议
        if (modTunnel.enable && tunnel2.adapter) {
            const modDevices = await mongo.db.collection("device").find({
                tunnel_id: modTunnel._id,
                enable: true
            }).toArray();

            if (modDevices.length) {
                //启动相关设备
                modDevices.forEach(d => device.create(tunnel2, d));

                //启动相关项目
                const projects = await mongo.db.collection("project").find({
                    enable: true,
                    devices: {
                        $elemMatch: {
                            device_id: {
                                $in: modDevices.map(d => {
                                    return d._id
                                })
                            }
                        }
                    }
                }).toArray()
                projects.forEach(prj => {
                    //检查项目
                    if (!project.get(prj._id))
                        project.create(prj)
                });
            }
        }

        tunnel2.on('close', () => {
            if (tunnel2.model) {
                createEvent({tunnel_id: tunnel2.model._id, event: "下线"})
                //exports.removeTunnel(tunnel2.model._id)

                delete tunnels[tunnel2.model._id];
            }
        })
    });

    return acceptor;
}



