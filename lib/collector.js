const EventEmitter = require("events");
const cron = require('./cron');
const interval = require('./interval');

class Collector extends EventEmitter {
    adapter;

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

    constructor(adapter, options) {
        super();
        this.adapter = adapter;
        Object.assign(this.options, options)

        this.start();
    }

    start() {
        if (this.cronHandle)
            this.cronHandle.cancel();

        if (this.options.interval)
            this.cronHandle = interval.check(this.options.interval * 1000, () => this.execute());
        else if (this.options.crontab)
            this.cronHandle = cron.schedule(this.options.crontab, () => this.execute());
        else
            throw new Error("间隔 和 定时 至少要有一项")
    }

    execute() {
        if (this.reading) {
            log.trace(this.options, '正在执行采集');
            return;
        }
        this.reading = true;
        this.adapter.read(this.options.code, this.options.address, this.options.length).then(data => {
            //console.log('collector data', data)
            this.emit("data", data);
        }).catch(err => {
            //console.log(err)
            this.emit("error", err);
        }).finally(() => {
            this.reading = false;
        });
    }

    cancel() {
        //log.trace(this.options, '取消采集');
        if (this.cronHandle)
            this.cronHandle.cancel();
        this.cronHandle = undefined;
    }
}


exports.create = function (adapter, options) {
    return new Collector(adapter, options)
}