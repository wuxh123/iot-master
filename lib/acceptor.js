const fs = require("fs");
const path = require("path");

/**
 * 创建接收器（服务）
 * @param type 类型
 * @param options 参数
 * @return {Acceptor}
 */
exports.create = function (type, options) {
    //检查js脚本是否存在
    const mod = path.join('../acceptors', type + '.js');
    if (!fs.existsSync(mod)) {
        throw new Error("不支持的接收器类型：" + type);
    }

    return require(mod)(options);
}

exports.recovery = function () {
    const files = fs.readdirSync(path.join(global.data_path, "acceptors"))
    files.forEach(cfg => {
        const acceptor = exports.create(cfg.type, cfg.options)
        acceptor.on('connect', tunnel => {

            tunnel.on('register', sn => {
                //TODO 创建设备，存入数据库

                //TODO 创建collector

            })
        })


    })
}



