const EventEmitter = require("events");
const cron = require('./cron');

class Collector extends EventEmitter {
    adapter;
    options = {
        crontab: '*/5 * * * * *',
        slave: 1,
        address: '',
        length: 0,
        values: {}
    };

    reading = false;

    cronHandle;

    constructor(adapter, options) {
        super();
        this.adapter = adapter;
        Object.assign(this.options, options)

        //scheduleCollector(crontab, this)
        this.cronHandle = cron.schedule(this.options.crontab, () => this.execute());
    }

    execute() {
        if (this.reading) return;
        this.reading = true;
        this.adapter.read(this.options.slave, this.options.address, this.options.length).then(data => {
            //解析结果，
            const values = {};
            this.options.values.forEach(v => {
                values[v.name] = data[v.offset];
                if (v.scale) values[v.name] *= v.scale;
            });
            this.emit("values", values);
        }).catch(err => {
            this.emit("error", err);
        }).finally(() => {
            this.reading = false;
        });
    }

    cancel() {
        this.cronHandle.cancel();
    }
}


exports.create = function (adapter, options) {
    return new Collector(adapter, options)
}