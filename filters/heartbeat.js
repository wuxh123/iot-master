module.exports = class HeartBeat {
    tunnel;
    options = {};

    last = Date.now();

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
        //支持转义字符，比如：换行 回车 十六进制
        if (options.text)
            this.options.text = eval("`" + options.text + "`");
    }

    handle(data, next) {
        let now = Date.now();
        if (now > this.options.interval + this.last) {
            if (this.options.text && this.options.text === data.toString() ||
                this.options.regex && this.options.regex.test(data.toString())
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