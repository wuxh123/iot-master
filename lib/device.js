const EventEmitter = require("events");


module.exports = class Device extends EventEmitter {
    /**
     * @type Adapter
     */
    adapter;

    options = {};

    collectors = [];
    values = {};


    constructor(adapter, options) {
        super();
        this.adapter = adapter;
        Object.assign(this.options, options);
    }

    execute(cmd, args) {
        const command = this.options.commands[cmd];
        this.adapter.write(this.slave, command.address, command.value)
        //(...args);
    }
}