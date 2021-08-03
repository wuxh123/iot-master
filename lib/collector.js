const EventEmitter = require("events");
const cron = require('./cron');
const interval = require('./interval');

class Collector extends EventEmitter {
    adapter;
    options = {
        interval: 30,
        crontab: '*/5 * * * * *',
        slave: 1,
        code: 1,
        address: 0,
        length: 0
    };

    reading = false;

    cronHandle;

    constructor(adapter, options) {
        super();
        this.adapter = adapter;
        Object.assign(this.options, options)

        //scheduleCollector(crontab, this)
        if (this.options.interval)
            this.cronHandle = interval.check(this.options.interval, () => this.execute());
        else if (this.options.crontab)
            this.cronHandle = cron.schedule(this.options.crontab, () => this.execute());
    }

    execute() {
        if (this.reading) return;
        this.reading = true;
        this.adapter.read(this.options.slave, this.options.code, this.options.address, this.options.length).then(data => {
            //TODO 解析结果，根据设备对应元件的变量映射表
            this.emit("values", data);
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