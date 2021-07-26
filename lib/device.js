const EventEmitter = require("events");
const Collector = require("./collector");

module.exports = class Device extends EventEmitter {
    /**
     * @type Adapter
     */
    adapter;

    options = {};

    collectors = [];

    values = {
        a: {
            value: 1,
            last: Date.now(),
            watchers: []
        }
    };


    constructor(adapter, options) {
        super();
        this.adapter = adapter;
        Object.assign(this.options, options);

        this.options.collectors.forEach(c => {
            const collector = new Collector(adapter, c);
            collector.on("values", values => {
                //Object.assign(this.values, values);

                //向project汇报
                this.emit("values", values);

                //执行watchers

            })
        })

    }

    execute(cmd, param) {
        const command = this.options.commands[cmd];
        this.adapter.write(this.options.slave, command.address, param || command.value);
    }



}