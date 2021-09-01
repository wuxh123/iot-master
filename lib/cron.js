const schedule = require("node-schedule");

const schedules = {};

exports.schedule = function (crontab, callback) {
    //检查定时器
    if (!schedules.hasOwnProperty(crontab)) {
        schedules[crontab] = {
            callbacks: [],
            job: schedule.scheduleJob(crontab, () => {
                log.trace({crontab}, '定时器')
                schedules[crontab].callbacks.forEach(c => {
                    try {
                        c()
                    } catch (e) {
                        log.error(e.message)
                    }
                })

            }),
        };
    }
    schedules[crontab].callbacks.push(callback);

    return {
        cancel: function () {
            const callbacks = schedules[crontab].callbacks;
            const index = callbacks.indexOf(callback);
            if (index > -1) callbacks.splice(index, 1);
            if (callbacks.length === 0) {
                schedules[crontab].job.cancel();
                delete schedules[crontab];
            }
        }
    }
}

exports.scheduleTime = function (minutes, callback) {
    if (minutes < 0) minutes = 0;
    else if (minutes > 1439) minutes = 1439;
    const m = Math.floor(minutes % 60);
    const h = Math.floor(minutes / 60);
    const c = `${m} ${h} * * *`;
    return exports.schedule(c, callback);
}