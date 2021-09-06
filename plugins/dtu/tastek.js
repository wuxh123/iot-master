class Tastek {
    tunnel;
    options = {
        prefix: '@DTU:0000:', //有人模块支持网络AT，usr.cn#AT\r\n
    };

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
    }

    query(cmd) {
        let command;
        switch (cmd) {
            case 'firmware':
                command = "+CGMR";
                break;
            case 'imei':
                command = "+GSN";
                break;
            case 'iccid':
                command = "+ICCID";
                break;
            case 'gps':
                command = "+GPSINFO";
                break;
            case 'rssi':
                command = "+CSQ";
                break;
            default:
                return;
        }

        this.tunnel.write(this.options.prefix + command + '\r\n');
    }

    handle(data, next) {
        const text = data.toString();
        if (text.startsWith(this.options.prefix)) {
            const result = text.substring(this.options.prefix.length);
            let results = result.split('\r\n');
            if (results.length < 2) return;
            results = results[1].split(':');
            switch (results[0]) {
                case '+CGMR':
                    this.tunnel.emit('control', {'firmware': results[1]});
                    break;
                case '+GSN':
                    this.tunnel.emit('control', {'imei': results[1]});
                    break;
                case '+ICCID':
                    this.tunnel.emit('control', {'iccid': results[1]});
                    break;
                case '+GPSINFO'://
                    let longitude, latitude;
                    const strs = results[1].split(',')
                    this.tunnel.emit('control', {'gps': [Number(strs[0]), Number(strs[1])]});
                    break;
                case '+CSQ':
                    this.tunnel.emit('control', {'rssi': parseInt(results[1])});
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
    return new Tastek(tunnel, options);
}