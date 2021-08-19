class Config {
    tunnel;
    options = {
        prefix: 'config,', //Luat模块支持远程配置命令 config,get,imei\r\n
        suffix: '\r\n',
    };

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
    }

    query(cmd) {
        let command;
        switch (cmd) {
            case 'firmware':
                command = "firmwarever";
                break;
            case 'imei':
                command = "imei";
                break;
            case 'iccid':
                command = "iccid";
                break;
            case 'lbs':
                command = "lbsloc";
                break;
            case 'rssi':
                command = "csq";
                break;
            default:
                return;
        }

        this.tunnel.write(this.options.prefix + 'get,' + command + '\r\n');
    }

    handle(data, next) {
        const text = data.toString();
        if (text.startsWith(this.options.prefix) && text.endsWith(this.options.suffix)) {
            //this.tunnel.emit('config', text.substring(this.options.prefix.length, -this.options.suffix.length));
            const result =  text.substring(this.options.prefix.length, -this.options.suffix.length);
            const results = result.split(',')
            switch (results[0]) {
                case 'imei':
                    if (results[1] === 'ok')
                        this.tunnel.emit('imei', results[2]);
                    break;
                case 'iccid':
                    if (results[1] === 'ok')
                        this.tunnel.emit('iccid', results[2]);
                    break;
                case 'lbsloc':
                    if (results[1] === 'ok')
                        this.tunnel.emit('gps', {longitude: Number(results[2]), latitude: Number(results[3])});
                    break;
                case 'firmwarever':
                    if (results[1] === 'ok')
                        this.tunnel.emit('firmware', results[2]);
                    break;
                case 'csq':
                    if (results[1] === 'ok')
                        this.tunnel.emit('rssi', Number(results[2]));
                    break;
                default:
                    //DO nothing
                    break;
            }

            return;
        }

        next(data)
    }
}

module.exports = function (tunnel, options) {
    return new Config(tunnel, options);
}