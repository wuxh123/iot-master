/**
 *
 * @type {Interval[]}
 */
const internals = [];

class Interval {
    interval;
    callback;

    last;

    constructor(interval, callback) {
        this.interval = interval;
        this.callback = callback;
        this.last = Date.now();
    }

    cancel() {
        for (let i = 0; i < internals.length; i++) {
            if (internals[i] === this) {
                internals.splice(i, 1)
                return;
            }
        }
    }
}

setInterval(() => {
    const now = Date.now();

    for (let i = 0; i < internals.length; i++) {
        const int = internals[i];
        if (now - int.last > int.interval) {
            int.last = now;
            process.nextTick(int.callback)
        }
    }
}, 1000);

/**
 *
 * @param {number} int
 * @param {function} cb
 * @return Interval
 */
exports.register = function (int, cb) {
    const interval = int * 1000;
    const i = new Interval(interval, cb);
    i.last += Math.floor(interval * Math.random()); //增加随机性，避免定时器过于集中
    internals.push(i);
    return i;
}

exports.intervals = internals;