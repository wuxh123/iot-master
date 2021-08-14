const fs = require("fs");
const path = require("path");
const EventEmitter = require("events");

const device = require("./device");
const project = require("./project");

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

    acceptor.on('connect', tunnel => {
        tunnel.on('register', async sn => {
            //TODO 通道暂不支持自定义配置（过滤器，协议等）

            const ret = await mongo.db.collection("tunnel").findOneAndUpdate({
                acceptor_id: model._id,
                sn
            }, {$set: {
                remote: tunnel.conn.remoteAddress,
                last: new Date()
            }}, {
                upsert: true,
                returnDocument: 'after'
            });

            const modTunnel = ret.value;
            tunnel.model = modTunnel;

            //添加上线记录
            mongo.db.collection("event").insertOne({tunnel_id: modTunnel._id, event: "上线"}).then()

            //首次上线
            if (!modTunnel.hasOwnProperty('enable')) {
                //补充启用
                modTunnel.enable = true;
                mongo.db.collection("tunnel").updateOne({_id: modTunnel._id}, {$set: {enable: true}}).then();

                //创建设备，存入数据库
                if (tunnel.adapter && model.devices && model.devices.length) {
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
            tunnels[modTunnel._id] = tunnel

            //TODO 此处解析通道配置，初始化通道协议
            if (modTunnel.enable && tunnel.adapter) {
                const modDevices = await mongo.db.collection("device").find({
                    tunnel_id: modTunnel._id,
                    enable: true
                }).toArray();

                if (modDevices.length) {
                    //启动相关设备
                    modDevices.forEach(d => device.create(tunnel, d));

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
        });

        tunnel.on('close', () => {
            if (tunnel.model) {
                mongo.db.collection("event").insertOne({tunnel_id: tunnel.model._id, event: "下线"}).then()
                //exports.removeTunnel(tunnel.model._id)
                delete tunnels[tunnel.model._id];
            }
        })
    });

    return acceptor;
}



