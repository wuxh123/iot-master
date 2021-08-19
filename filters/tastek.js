class UsrCn {
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
            case 'sn':
                command = "+DEVICEID";
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
            //this.tunnel.emit('at', text.substring(this.options.prefix.length, -this.options.suffix.length));
            const result = text.substring(this.options.prefix.length);
            let results = result.split('\r\n');
            results = results[0].split(':');
            switch (results[0]) {
                case '+CGMR':
                    this.tunnel.emit('firmware', results[1]);
                    break;
                case '+DEVICEID':
                    this.tunnel.emit('sn', results[1]);
                    break;
                case '+GSN':
                    this.tunnel.emit('imei', results[1]);
                    break;
                case '+ICCID':
                    this.tunnel.emit('iccid', results[1]);
                    break;
                case '+GPSINFO'://
                    let longitude, latitude;
                    const strs = results[1].split(',')
                    this.tunnel.emit('gps', {longitude: Number(strs[0]), latitude: Number(strs[1])});
                    break;
                case '+CSQ':
                    this.tunnel.emit('rssi', parseInt(results[1]));
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
    return new UsrCn(tunnel, options);
}