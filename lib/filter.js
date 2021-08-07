const path = require("path");
const fs = require("fs");

class Filter {
    /**
     * @type Tunnel
     */
    tunnel;

    /**
     * 处理数据
     * @param {Buffer|string} data
     * @param {function} next
     */
    handle(data, next) {
    }
}

exports.create = function (tunnel, type, options) {
    //检查js文件
    //return new require('../acceptors/filters/' + option.type)(tunnel, option.options);
    const mod = path.join(__dirname, '..', 'filters', type + '.js');
    if (!fs.existsSync(mod)) {
        throw new Error("不支持的接收器类型：" + type);
    }
    return require(mod)(tunnel, options);
}