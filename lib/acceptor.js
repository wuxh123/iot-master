const fs = require("fs");
const path = require("path");
const Tunnel = require("./tunnel");
const EventEmitter = require("events");

const mongo = require_plugin("mongo");

class Acceptor extends EventEmitter{

}

/**
 * 创建接收器（服务）
 * @param model 参数
 * @return {Acceptor}
 */
exports.create = function (model) {
    //检查js脚本是否存在
    const mod = path.join('../acceptors', model.type + '.js');
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
                if (model.devices && model.devices.length) {
                    //依次创建设备
                    const modDevices = model.devices.map(d=>{
                        return {
                            tunnel_id: d.tunnel_id,
                            slave: d.slave,
                            enable: true,
                        };
                    })
                    await mongo.db.collection("device").insertMany(modDevices);
                }

            }

            if (modTunnel.enable) {
                const modDevices = await mongo.db.collection("device").find({tunnel_id: modTunnel._id, enable: true}).toArray();
                //TODO 依次创建设备


            }
        });
    });

    return acceptor;
}



