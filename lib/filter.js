
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

exports.create = function (tunnel, option) {
    //TODO 检查js文件
    return new require('../acceptors/filters/' + option.type)(tunnel, option.options);
}