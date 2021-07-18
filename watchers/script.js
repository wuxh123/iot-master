const vm = require('vm');

module.exports = class ScriptWatcher {

    constructor(values, options) {
        vm.runInNewContext(options.script, {a:1, b:2})
    }

}