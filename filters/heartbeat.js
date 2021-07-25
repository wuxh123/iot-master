module.exports = class HeartBeat {
    tunnel;
    options = {};

    last = Date.now();

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
    }

    handle(data, next) {
        let now = Date.now();
        if (now > this.options.interval + this.last) {
            if (this.options.text && this.options.text === data.toString() ||
                this.options.regex && this.options.regex.test(data.toString()) ||
                this.options.hex && this.options.hex.toLowerCase() === data.toString('hex')
            ) {
                this.last = now;
                this.tunnel.emit('heartbeat', data);
                return;
            }
        }
        this.last = now;

        next(data)
    }
}