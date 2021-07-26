const cache = {};

//开放给全局使用
global.plugins = cache;

//声明引入插件接口，方便自定义脚本使用
global.require_plugin = function (id) {
    //优先加载缓存
    if (cache.hasOwnProperty(id))
        return cache[id];

    //加载插件
    let plugin = require("../plugins/" + id);

    //函数形式
    if (typeof plugin === 'function') {
        //加载配置文件
        let cfg = {};
        try {
            cfg = require("../data/plugins/" + id);
        } catch (e) {
            //没有配置文件
        }
        plugin = plugin(cfg);
    }
    return plugin;
}