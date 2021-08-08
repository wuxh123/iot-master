const fs = require("fs");
const path = require("path");
const EventEmitter = require("events");

const device = require("./device");
const project = require("./project");

const mongo = require_plugin("mongodb");

class Acceptor extends EventEmitter {

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

    acceptor.on('connect', tunnel => {
        tunnel.on('register', async sn => {

            //TODO 通道暂不支持自定义配置（过滤器，协议等）

            let modTunnel = await mongo.db.collection("tunnel").findOne({acceptor_id: model._id, sn});
            if (!modTunnel) {
                const res = await mongo.db.collection("tunnel").insertOne({acceptor_id: model._id, sn, enable: true});
                modTunnel = await mongo.db.collection("tunnel").findOne({_id: res.insertedId});
                //TODO 可以使用 findOneAndUpdate, upsert功能

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
    });

    return acceptor;
}



