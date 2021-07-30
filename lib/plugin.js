const path = require("path");

const cache = {};

//开放给全局使用
global.plugins = {};

//声明引入插件接口，方便自定义脚本使用
global.require_plugin = function (id) {
    //优先加载缓存
    if (global.plugins.hasOwnProperty(id))
        return global.plugins[id];

    //加载插件
    let plugin = global.plugins[id] = require("../plugins/" + id);

    //函数形式
    if (typeof plugin.config === 'function') {
        //加载配置文件
        let cfg = {};
        try {
            cfg = require(path.join(global.data_path, id + '.config'));
        } catch (e) {
            //没有配置文件
        }
        plugin.config(cfg);
    }
    return plugin;
}