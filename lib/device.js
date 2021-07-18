const EventEmitter = require("events");


module.exports = class Device extends EventEmitter {
    adapter;
    options = {};

    collectors = [];
    values = {};

    commands = {}

    constructor(adapter, options) {
        super();
        this.adapter = adapter;
        Object.assign(this.options, options);
    }

    execute(command, args) {
        this.commands[command](...args);
    }
}