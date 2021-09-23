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
    const i = new Interval(int * 1000, cb)
    internals.push(i);
    return i;
}

exports.intervals = internals;