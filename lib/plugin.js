const fs = require("fs");
const path = require("path");
const semver = require("semver");
const child_process = require("child_process");
const _ = require("lodash");

//声明引入插件接口，方便自定义脚本使用
global.require_plugin = function (id) {
    return require("../plugins/" + id);
}

global.load_config = function (name) {
    let cfg = {};
    try {
        cfg = require(path.join(global.data_path, name + '.config'));
    } catch (e) {
        //没有配置文件
    }
    return cfg;
}


function loadManifest(dir) {
    const spec = path.join(dir, 'package.json');
    const stat = fs.statSync(spec, {throwIfNoEntry: false});
    if (stat && stat.isFile()) {
        const json = fs.readFileSync(spec, 'utf8');
        return JSON.parse(json);
    }

    return null
}

function storeManifest(dir, m) {
    const data = JSON.stringify(m, null, '\t');
    const spec = path.join(dir, 'package.json');
    fs.writeFileSync(spec, data);
}

const plugins_dir = "plugins";

/**
 * 安装插件
 * @param {string} p 插件名（目录）
 */
exports.install = function (p) {
    const manifest = loadManifest('.');
    const plugin = loadManifest(path.join(plugins_dir, p));
    //manifest.plugins[p] = plugin.version || '1.0.0';

    //if (!plugin) throw new Error("找不到说明文件");
    if (!plugin) return;

    let newPack = false;
    Object.keys(plugin.dependencies).forEach(pack => {
        const ver = plugin.dependencies[pack];
        const mVer = manifest.dependencies[pack];
        if (!mVer || semver.ltr(semver.minVersion(mVer), ver)) {
            manifest.dependencies[pack] = ver;
            newPack = true;
        }
    });

    //安装依赖包
    if (newPack) {
        storeManifest('.', manifest);
        //child_process.execSync("npm install");
        let ls = child_process.spawn('npm', ['install'], {
            stdio: ['pipe', 'pipe', 'pipe']
        });
        ls.stdout.on("data", console.log)
        ls.stderr.on("data", console.error)
    }
}

function clear() {
    exports.plugins = [];
    exports.acceptors = [];
    exports.filters = [];
    exports.protocols = [];
    exports.databases = [];
    exports.memoryDatabases = [];
    exports.historyDatabases = [];
    exports.smsPushers = [];
    exports.voicePushers = [];
    exports.webAPIs = [];
}

exports.createService = function create(type, name) {
    const p = exports[type + 's'].find(p => p.name === name);
    if (!p) throw new Error("未找到组件：" + name);
    if (!fs.existsSync(p.script)) throw new Error("组件脚本缺失：" + p.script);
    const args = [...arguments];
    return new (require(path.join('..', p.script)))(... args.slice(2));
}

exports.createAcceptor = function (name, options) {
    return exports.createService('acceptor', name, options);
}

exports.createFilter = function (name, tunnel, options) {
    return exports.createService('filter', name, tunnel, options);
}

exports.createProtocol = function (name, tunnel, options) {
    return exports.createService('protocol', name, tunnel, options);
}

exports.createDatabase = function (name, options) {
    return exports.createService('database', name, options);
}

exports.createHistoryDatabase = function (name, options) {
    return exports.createService('historyDatabase', name, options);
}

exports.createMemoryDatabase = function (name, options) {
    return exports.createService('memoryDatabase', name, options);
}

exports.createSMSPusher = function (name, options) {
    return exports.createService('smsPusher', name, options);
}

exports.createVoicePusher = function (name, options) {
    return exports.createService('voicePusher', name, options);
}


/**
 * 加载插件
 * @param {string} p 插件名（目录）
 */
exports.load = function (p) {
    const spec = path.join(plugins_dir, p, 'package.json');
    const stat = fs.statSync(spec, {throwIfNoEntry: false});
    if (stat && stat.isFile()) {
        const json = fs.readFileSync(spec, 'utf8');
        const manifest = JSON.parse(json);
        exports.plugins.push(manifest);

        //加载配置，服务，协议，接口，等
        function append(collection, array) {
            array && array.forEach(a => {
                const v = _.cloneDeep(a);
                v.script = path.join(plugins_dir, p, a.script);
                collection.push(v);
            })
        }

        append(exports.acceptors, manifest.acceptors);
        append(exports.filters, manifest.filters);
        append(exports.protocols, manifest.protocols);
        append(exports.databases, manifest.databases);
        append(exports.historyDatabases, manifest.historyDatabases);
        append(exports.memoryDatabases, manifest.memoryDatabases);
        append(exports.smsPushers, manifest.smsPushers);
        append(exports.voicePushers, manifest.voicePushers);
        append(exports.webAPIs, manifest.webAPIs);
    } else {
        exports.plugins.push({name: p})
    }
}

/**
 * 加载插件
 */
exports.loadAll = function () {
    clear();

    const folders = fs.readdirSync(plugins_dir);
    folders.forEach(folder => {
        const filePath = path.join(plugins_dir, folder);
        let stat = fs.statSync(filePath);
        if (stat.isDirectory()) {
            //console.log("[plugin]", folder);
            exports.load(folder);
            //exports.install(folder); //调试用
        }
    })
}

exports.loadAll();