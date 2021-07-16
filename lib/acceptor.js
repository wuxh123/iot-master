
exports.create = function(type, options) {
    //从acceptors中找到脚本
    //TODO 检查js脚本是否存在
    return new require('../acceptors/' + type)(options);
}

