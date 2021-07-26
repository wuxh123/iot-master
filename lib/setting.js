const fs = require("fs");
const path = require("path");

const filename = path.join(global.data_path, "config.json");

const defaultOptions = {};

global.setting = defaultOptions;


exports.save = function (options) {
    fs.writeFileSync(filename, JSON.stringify(options, null, '\t'));
}

exports.reset = function () {
    global.setting = defaultOptions;
    exports.save(defaultOptions)
}

//加载配置
try {
    const text = fs.readFileSync(filename, "utf-8");
    const data = JSON.parse(text);
    global.setting = Object.assign(defaultOptions, data);
} catch (e) {
    exports.reset();
}