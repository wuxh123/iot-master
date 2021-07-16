const schedule = require("node-schedule");
const EventEmitter = require("events");

const schedules = {};

function scheduleCollector(crontab, collector) {
    //检查定时器
    if (!schedules.hasOwnProperty(crontab)) {
        schedules[crontab] = {
            collectors: [],
            job: schedule.scheduleJob(crontab, () => {
                schedules[crontab].collectors.forEach(c => c.execute())
            }),
        };
    }
    schedules[crontab].collectors.push(collector);
}


class Collector extends EventEmitter {
    adapter;
    crontab = '*/5 * * * * *';
    command = '';
    values = {};

    reading = false;

    constructor(adapter, crontab, command, values) {
        super();

        this.adapter = adapter;
        this.crontab = crontab;
        this.command = command;
        this.values = values;
        //scheduleCollector(crontab, this)
    }

    execute() {
        if (this.reading) return;
        this.reading = true;
        this.adapter.read(this.command).then(data => {
            //解析结果，
            const values = {};
            this.values.forEach(v => {
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
        const collectors = schedules[this.crontab].collectors;
        const index = collectors.indexOf(this);
        if (index > -1) collectors.splice(index, 1);
        if (collectors.length === 0) {
            schedules[this.crontab].job.cancel();
            delete schedules[this.crontab];
        }
    }
}


//options: interval, command(slave, code, data, length), values[offset, type(word:uint16, bit), ]
exports.create = function (adapter, crontab, command, values) {
    const collector = new Collector(adapter, crontab, command, values);
    scheduleCollector(crontab, collector);
    return collector;
}
