const fs = require("fs");
const path = require("path")

/**
 * 创建接收器（服务）
 * @param type 类型
 * @param options 参数
 * @return {Acceptor}
 */
exports.create = function(type, options) {
    //TODO 检查js脚本是否存在
    return new require('../acceptors/' + type)(options);
    //TODO 监听 connect 得到通道
    //TODO 监听 tunnel.register 得到注册码，同步数据库


}

exports.recovery = function () {
    const files = fs.readdirSync(path.join(global.data_path, "acceptors"))
    files.forEach(f=>{
        exports.create(f.type, f.options)
    })
}



