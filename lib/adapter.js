const fs = require("fs");
const path = require("path");
const EventEmitter = require("events");

const protocols_dir = "protocols";
exports.protocols = [];

/**
 * 加载协议
 */
exports.load = function () {
    exports.protocols = [];
    const folders = fs.readdirSync(protocols_dir);
    folders.forEach(folder => {
        const filePath = path.join(protocols_dir, folder);
        let stat = fs.statSync(filePath);
        if (stat.isDirectory()) {
            const spec = path.join(filePath, 'package.json');
            stat = fs.statSync(spec, {throwIfNoEntry: false});
            if (stat && stat.isFile()) {
                console.log("[protocol]", folder);
                const json = fs.readFileSync(spec, 'utf8');
                const data = JSON.parse(json);
                data.protocols && data.protocols.forEach(p => {
                    p.codes = p.codes || data.codes;
                    p.version = p.version || data.version;
                    p.script = path.join(folder, p.script);
                    exports.protocols.push(p);
                    exports.protocols[p.name] = p;
                });
            }
        }
    })
}

class Adapter extends EventEmitter{

}

/**
 * 创建适配器
 * @param {Tunnel} tunnel
 * @param {Object} options 协议配置
 * @return {Adapter}
 */
exports.create = function (tunnel, options) {
    //TODO 检查js文件
    return require('../protocols/' + options.type)(this, options.options);
}

//加载协议
exports.load();
//console.log(exports.protocols)