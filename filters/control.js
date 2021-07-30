module.exports = class Control {
    tunnel;
    options = {
        prefix: 'user.cn#', //Luat模块支持远程配置命令 config,get,imei\r\n 有人模块支持网络AT，usr.cn#AT\r\n
        suffix: '\r\n'
    };

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
        //支持转义字符，比如：换行 回车 十六进制
        if (options.prefix)
            this.options.prefix = eval("`" + options.prefix + "`");
        if (options.suffix)
            this.options.suffix = eval("`" + options.suffix + "`");
    }

    handle(data, next) {
        const text = data.toString();
        if (text.startsWith(this.options.prefix) && text.endsWith(this.options.suffix)) {
            this.tunnel.emit('control', text.substring(this.options.prefix.length, -this.options.suffix.length));
            return;
        }

        next(data)
    }
}