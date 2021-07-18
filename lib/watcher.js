
class Watcher {

}


exports.create = function(type, options) {
    //从watchers中找到脚本
    //TODO 检查js脚本是否存在
    return new require('../watchers/' + type)(options);
}

