const EventEmitter = require("events");
const cron = require('./cron');
const pacemaker = require('./pacemaker');

module.exports = class Collector extends EventEmitter {
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

    /**
     * 开始
     */
    start() {
        if (this.cronHandle)
            this.cronHandle.cancel();

        if (this.options.interval)
            this.cronHandle = pacemaker.register(this.options.interval, () => this.execute());
        else if (this.options.crontab)
            this.cronHandle = cron.schedule(this.options.crontab, () => this.execute());
        else
            throw new Error("间隔 和 定时 至少要有一项")
    }

    /**
     * 执行
     */
    execute() {
        if (this.reading) {
            log.trace(this.options, 'collector execute');
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

    /**
     * 刷新
     * @returns {Promise}
     */
    refresh(){
        return this.adapter.read(this.options.code, this.options.address, this.options.length, true);
    }

    /**
     * 取消
     */
    cancel() {
        if (this.cronHandle)
            this.cronHandle.cancel();
        this.cronHandle = undefined;
    }
}
