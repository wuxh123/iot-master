class Acceptor {

    constructor(options) {
    }

    close() {

    }
}


/**
 * 创建接收器（服务）
 * @param type 类型
 * @param options 参数
 * @return {Acceptor}
 */
exports.create = function(type, options) {
    //TODO 检查js脚本是否存在
    return new require('../acceptors/' + type)(options);
}

