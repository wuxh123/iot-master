const EventEmitter = require("events");
const cron = require('./cron');
const interval = require('./interval');

class Collector extends EventEmitter {
    adapter;
    slave = 1;
    options = {
        interval: 0,
        crontab: '',
        slave: 1,
        code: 1,
        address: 0,
        length: 0
    };

    reading = false;

    cronHandle;

    constructor(adapter, slave, options) {
        super();
        this.adapter = adapter;
        this.slave = slave;
        Object.assign(this.options, options)

        //scheduleCollector(crontab, this)
        if (this.options.interval)
            this.cronHandle = interval.check(this.options.interval * 1000, () => this.execute());
        else if (this.options.crontab)
            this.cronHandle = cron.schedule(this.options.crontab, () => this.execute());
        else
            throw new Error("间隔 和 定时 至少要有一项")
    }

    execute() {
        if (this.reading) return;
        this.reading = true;
        this.adapter.read(this.options.slave, this.options.code, this.options.address, this.options.length).then(data => {
            this.emit("data", data);
        }).catch(err => {
            //console.log(err)
            this.emit("error", err);
        }).finally(() => {
            this.reading = false;
        });
    }

    cancel() {
        if (this.cronHandle)
            this.cronHandle.cancel();
    }
}


exports.create = function (adapter, slave, options) {
    return new Collector(adapter, slave, options)
}