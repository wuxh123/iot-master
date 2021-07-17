const schedule = require("node-schedule");

const schedules = {};

exports.schedule = function(crontab, callback) {
    //检查定时器
    if (!schedules.hasOwnProperty(crontab)) {
        schedules[crontab] = {
            callbacks: [],
            job: schedule.scheduleJob(crontab, () => {
                schedules[crontab].callbacks.forEach(c => c.execute())
            }),
        };
    }
    schedules[crontab].callbacks.push(callback);

    return {
        cancel: function () {
            const callbacks = schedules[crontab].callbacks;
            const index = callbacks.indexOf(this);
            if (index > -1) callbacks.splice(index, 1);
            if (callbacks.length === 0) {
                schedules[crontab].job.cancel();
                delete schedules[crontab];
            }
        }
    }
}
