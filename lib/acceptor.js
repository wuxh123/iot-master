const fs = require("fs");
const path = require("path");
const EventEmitter = require("events");

const tunnel = require("./tunnel");
const device = require("./device");
const project = require("./project");
const plugin = require("./plugin");

const mongo = require_plugin("mongodb");

class Acceptor extends EventEmitter {

}

const acceptors = {};

exports.getAcceptor = function (id) {
    return acceptors[id];
}

exports.removeAcceptor = function (id) {
    log.info({id}, 'removeAcceptor')
    const acc = acceptors[id];
    if (acc) {
        acc.close();
        delete acc[id];  //projects[id]=null
    }
}


/**
 * 创建接收器（服务）
 * @param model 参数
 * @return {Acceptor}
 */
exports.create = function (model) {
    log.info({id: model._id, type: model.type, port: model.port}, 'createAcceptor');

    const acceptor = plugin.createAcceptor(model.type, model);

    //保存到全局
    acceptors[model._id] = acceptor;

    acceptor.on('error', err => {
        log.error(err.message, "启动失败");
    });

    acceptor.on('connect', async (sn, socket) => {
        let newTunnel = false;

        //根据SN找通道记录
        let modTunnel = await mongo.db.collection("tunnel").findOne({acceptor_id: model._id, sn},);
        if (modTunnel) {
            if (!modTunnel.enable) {
                socket.write("denied sn")
                log.trace("tunnel denied")
                return;
            }
            //更新上线时间
            mongo.db.collection("tunnel").updateOne({_id: modTunnel._id}, {
                $set: {last: new Date(), remote: socket.remoteAddress, online: true}
            }).then().catch(err => this.emit('error', err));

        } else {
            modTunnel = {
                acceptor_id: model._id,
                sn,
                heartbeat: model.heartbeat,
                control: model.control,
                protocol: model.protocol,
                timeout: model.timeout,
                enable: true,
                last: new Date(),
                remote: socket.remoteAddress
            };
            const ret = await mongo.db.collection("tunnel").insertOne(modTunnel)
            modTunnel._id = ret.insertedId;

            newTunnel = true;
        }

        //打开通道
        const tnl = tunnel.open(socket, modTunnel);

        //初始化设备
        if (newTunnel) {
            //创建设备，存入数据库
            if (tnl.protocol && model.devices && model.devices.length) {
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


        //此处解析通道配置，初始化通道协议
        if (tnl.protocol) {
            const modDevices = await mongo.db.collection("device").find({
                tunnel_id: modTunnel._id,
                enable: true
            }).toArray();

            if (modDevices.length) {
                //启动相关设备
                modDevices.forEach(d => device.open(tnl, d));

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

    return acceptor;
}



