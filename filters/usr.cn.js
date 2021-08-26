class UsrCn {
    tunnel;
    options = {
        prefix: 'user.cn#', //有人模块支持网络AT，usr.cn#AT\r\n
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
                command = "AT+VER";
                break;
            case 'sn':
                command = "AT+SN";
                break;
            case 'imei':
                command = "AT+IMEI";
                break;
            case 'iccid':
                command = "AT+ICCID";
                break;
            case 'gps':
                command = "AT+LBS=1";
                break;
            case 'rssi':
                command = "AT+CSQ";
                break;
            case 'net':
                command = "AT+SYSINFO";
                break;
            default:
                return;
        }

        this.tunnel.write(this.options.prefix + command + '\r\n');
    }

    handle(data, next) {
        const text = data.toString();
        if (text.startsWith(this.options.prefix) && text.endsWith(this.options.suffix)) {
            //this.tunnel.emit('at', text.substring(this.options.prefix.length, -this.options.suffix.length));
            const result = text.substring(this.options.prefix.length, -this.options.suffix.length);
            let results = result.split('\r\n');
            results = results[0].split(':');
            switch (results[0]) {
                case '+VER': //+VER:V1.1.01.000000.0000
                    this.tunnel.emit('control', {'firmware': results[1]});
                    break;
                case '+IMEI': //+IMEI:864333040712457
                    this.tunnel.emit('control', {'imei': results[1]});
                    break;
                case '+ICCID'://+ICCID:8986003615195A571314
                    this.tunnel.emit('control', {'iccid': results[1]});
                    break;
                case '+LBS'://
                    let longitude, latitude;
                    results[1].split(',').forEach(str=>{
                        const strs = str.split('=')
                        if (strs[0] === 'LNG')
                            longitude = Number(strs[1]);
                        if (strs[0] === 'LAT')
                            latitude = Number(strs[1]);
                    })
                    if (longitude && latitude)
                        this.tunnel.emit('control', {'gps': [longitude, latitude]});
                    break;
                case '+SN'://+SN: 00402420011300024522
                    this.tunnel.emit('control', {'sn': results[1]});
                    break;
                case '+CSQ'://+CSQ: 27,99
                    this.tunnel.emit('control', {'rssi': parseInt(results[1])});
                    break;
                case '+SYSINFO'://+SYSINFO:4,LTE
                    this.tunnel.emit('control', {'net': results[1]});
                    break;
                case '+APN'://+APN:CMNET,,,0
                    this.tunnel.emit('control', {'apn': results[1]});
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