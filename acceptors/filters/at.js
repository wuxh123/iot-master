module.exports = class AT {
    tunnel;
    options = {
        prefix: 'usr.cn'
    };

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);

        //提前转成Buffer，提高效率
        //this.options.prefixBuffer = Buffer.from(this.options.prefix);
    }

    handle(data, next) {
        const text = data.toString();
        if (text.startsWith(this.options.prefix)) {
            this.tunnel.emit('at', text.substr(0, this.options.prefix.length));
            return;
        }
        // if (this.options.prefixBuffer.compare(data, 0, this.options.prefixBuffer.length) === 0) {
        //     this.tunnel.emit('at', data.slice(this.options.prefixBuffer.length).toString());
        //     return;
        // }

        next(data)
    }
}