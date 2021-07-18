const EventEmitter = require("events");


module.exports = class Device extends EventEmitter {
    /**
     * @type Adapter
     */
    adapter;

    options = {};

    collectors = [];
    values = {};

    commands = {
        close: {
            name: '关闭',
            address: 'BO002',
            value: false
        },
        closeAll: {
            name: '全部关闭',
            address: 'BO002',
            values: [false, false, false, false]
        }
    }

    constructor(adapter, options) {
        super();
        this.adapter = adapter;
        Object.assign(this.options, options);
    }

    execute(command, args) {
        this.commands[command](...args);
    }
}