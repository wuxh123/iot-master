const EventEmitter = require("events");
const interval = require('./interval');
const script = require("./script");

class Validator extends EventEmitter {
    model;

    //context
    variables;

    script;

    start = 0;
    delay = 0;

    reportAt = 0;
    reported = false;

    resetInterval;
    resetTimes = 0;

    checker; //定时检查

    constructor(model, variables) {
        super();
        this.model = model;
        this.variables = variables;
        
        this.script = script.compile(model.expression);
        this.start = 0;
        this.delay = model.delay * 1000; //换算成s，避免后续冗余的乘法计算
        this.resetInterval = model.resetInterval * 1000;
        this.resetTimes = 0;
    }

    execute() {
        try {
            const ctx = Object.assign({}, this.variables, this.variables.values())
            const ret = this.script.runInNewContext(ctx);
            if (ret) {
                this.start = 0;//去掉发生时间，重置延时
                this.reported = false;//去掉已经上报标识
                return;
            }
        } catch (err) {
            //this.error = err.message;
            log.error(err.message)
            this.emit('error', err);
            return;
        }
        //以下是不合法处理逻辑

        //延时处理，发生时间，当前时间
        const now = Date.now() * 0.001; //转换成秒
        if (this.delay) {
            if (!this.start) {
                this.start = now;
                return;
            }
            if (this.start + this.delay < now)
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
                if (this.reportAt + this.resetInterval < now)
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

        //TODO content 也可以使用表达式（脚本），但是变量要先拷贝，避免被覆盖
        this.emit('alarm', {
            name: this.model.name,
            content: this.model.content,
            level: this.model.level,
        });
    }

    cancel() {

    }

}