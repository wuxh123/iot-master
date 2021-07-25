module.exports = class LuatConfig {
    tunnel;
    options = {
        prefix: 'config,', //Luat模块支持远程配置命令 config,get,imei\r\n 有人模块支持网络AT，usr.cn#AT\r\n
        suffix: '\r\n'
    };

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
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