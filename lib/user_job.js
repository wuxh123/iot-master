const EventEmitter = require("events");
const cron = require("./cron");
const script = require("./script");
const _ = require("lodash");

module.exports = class UserJob extends EventEmitter {
    model;
    context;

    script;

    constructor(model, context) {
        super();

        this.model = model;
        this.context = context;

        //预编译参数数组
        if (model.params)
            this.script = script.compile('[' + model.params + ']')

        this.start();
    }

    start() {
        if (this.cronHandle)
            this.cronHandle.cancel();

        this.cronHandle = cron.schedule(this.model.crontab, () => {
            try {
                let params = [];
                if (this.script)
                    params = this.script.runInNewContext(this.context) //_.cloneDeep(this.context));

                //交给外面执行
                this.emit('execute', this.model.command, params);
            } catch (err) {
                this.emit('error', err);
            }
        })
    }

    cancel() {
        if (this.cronHandle)
            this.cronHandle.cancel();
        this.cronHandle = undefined;
    }
}