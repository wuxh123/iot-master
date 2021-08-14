
const vm = require("vm");

const scripts = {};

exports.compile = function (script) {
    return scripts[script] || (scripts[script] = new vm.Script(script));
}