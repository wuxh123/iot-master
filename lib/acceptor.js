const fs = require("fs");
const path = require("path");

/**
 * 创建接收器（服务）
 * @param options 参数
 * @return {Acceptor}
 */
exports.create = function (options) {
    //检查js脚本是否存在
    const mod = path.join('../acceptors', options.type + '.js');
    if (!fs.existsSync(mod)) {
        throw new Error("不支持的接收器类型：" + options.type);
    }

    const acceptor =  require(mod)(options);

    acceptor.on('connect', tunnel => {
        tunnel.on('register', sn => {

            //TODO 创建设备，存入数据库


        });
    });

    return acceptor;
}



