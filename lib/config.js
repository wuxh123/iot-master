const fs = require("fs");
const path = require("path");

const filename = path.join(global.data_path, "config.json");

const defaultOptions = {
    web: {
        enable: true,
        port: 8088,
        cors: true
    },
    traverser: {
        enable: true,
        port: 1843,
    }
};

global.config = defaultOptions;


exports.save = function (options) {
    fs.writeFileSync(filename, JSON.stringify(options || global.config, null, '\t'));
}

exports.reset = function () {
    global.config = defaultOptions;
    exports.save(defaultOptions)
}

//加载配置
try {
    const text = fs.readFileSync(filename, "utf-8");
    const data = JSON.parse(text);
    global.config = Object.assign({}, defaultOptions, data);
} catch (e) {
    exports.reset();
}