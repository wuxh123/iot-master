

/**
 * 创建适配器
 * @param {Tunnel} tunnel
 * @param {Object} options 协议配置
 * @return {Adapter}
 */
exports.create = function(tunnel, options) {
    //TODO 检查js文件
    return require('../adapters/' + options.type)(this, options.options);
}