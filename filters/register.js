class Register {
    tunnel;
    options = {};
    sn;

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
    }

    handle(data, next) {
        if (this.sn) {
            next(data);
            return;
        }

        let sn = data.toString();
        if (this.options.regex && !this.options.regex.test(sn)) {
            this.tunnel.write('invalid sn');
            //this.tunnel.close();
            return;
        }

        this.sn = sn;
        this.tunnel.emit('register', sn);
    }
}

module.exports = function (tunnel, options) {
    return new Register(tunnel, options);
}