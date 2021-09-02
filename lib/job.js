const EventEmitter = require("events");
const cron = require("./cron");

module.exports = class Job extends EventEmitter {
    model;
    context;

    script;

    constructor(model, context) {
        super();

        this.model = model;
        this.context = context;

        this.start();
    }

    start() {
        if (this.cronHandle)
            this.cronHandle.cancel();

    	this.cronHandle = cron.schedule(this.model.crontab, () => {
            try {
                this.script.runInNewContext(this.context);
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