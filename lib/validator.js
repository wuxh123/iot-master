const EventEmitter = require("events");
const cron = require('./cron');
const pacemaker = require('./pacemaker');
const script = require("./script");
const _ = require("lodash");

module.exports = class Validator extends EventEmitter {
    model;

    //context
    context;

    script;

    start = 0;

    reportAt = 0;
    reported = false;

    resetTimes = 0;

    checker; //定时检查

    constructor(model, context) {
        super();
        this.model = model;
        this.context = context;

        this.script = script.compile(model.expression);
        this.start = 0;
        this.resetTimes = 0;

        if (this.model.interval)
            this.cronHandle = pacemaker.register(this.model.interval, () => this.execute());
        else if (this.model.crontab)
            this.cronHandle = cron.schedule(this.model.crontab, () => this.execute());
    }

    execute() {
        try {
            //const ctx = _.cloneDeep(this.context) //Object.assign({}, this.variables, this.variables.values())
            const ret = this.script.runInNewContext(this.context);
            if (ret) {
                this.start = 0;//去掉发生时间，重置延时
                this.reported = false;//去掉已经上报标识
                return;
            }
        } catch (err) {
            //this.error = err.message;
            //log.error(err.message) 日志太多
            this.emit('error', err);
            return;
        }
        //以下是不合法处理逻辑

        //延时处理，发生时间，当前时间
        const now = Date.now() * 0.001; //转换成秒
        if (this.model.delay) {
            if (!this.start) {
                this.start = now;
                return;
            }
            if (this.start + this.model.delay > now)
                return;
        }

        //已经上报，则不再上报
        if (this.reported) {
            //重置逻辑
            if (this.model.resetInterval && this.model.resetTimes) {
                //如果已经超出重置次数，则不再执行
                if (this.resetTimes > this.model.resetTimes)
                    return;
                //如果还没到重置时间，则不提醒
                if (this.reportAt + this.model.resetInterval > now)
                    return;
                //重置
                this.start = now;
                this.resetTimes++
                return; //下次再执行，因为可能会有delay
            }
            return;
        }
        this.reported = true;
        this.reportAt = now

        //content 也可以含变量的字符串，但是变量要先拷贝，避免被覆盖
        this.emit('alarm', {
            name: this.model.name,
            content: this.model.content,
            level: this.model.level,
        });
    }

    cancel() {
        if (this.cronHandle)
            this.cronHandle.cancel();
        this.cronHandle = undefined;
    }

}