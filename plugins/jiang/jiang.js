const timeout = require('../../lib/interval');

function calcSum(buf, len) {
    let sum = 0;
    for (let i = 0; i < len; i++)
        sum += buf.readUInt8(i);
    return sum;
}

module.exports = class Jiang {

    tunnel;

    options = {
        concurrency: {
            enable: true,
            max: 50
        },
        timeout: 5000,
        transaction: {
            min: 0x01,
            max: 0xff
        }
    };

    transactionId = 0;

    _handlers = {};

    _queue = [];

    _doing = 0;

    _checker;

    onTunnelData = data => {
        this._handle(data)
    }

    onTunnelClose = () => {
        this._checker.cancel();
        this._checker = undefined;
        this.reset();
    }

    onTunnelOnline = () => {
        this.reset();
    }

    onTunnelOffline = () => {
        this.reset();
    }

    //清空所有请求，进行状态重置
    reset() {
        this._queue = [];
        for (let i in this._handlers) {
            if (this._handlers.hasOwnProperty(i)) {
                const h = this._handlers[i];
                this.reject(h.id, new Error("通道离线"));
            }
        }
        this._doing = 0;
    }

    constructor(tunnel, options) {
        Object.assign(this.options, options);
        this.open(tunnel)
    }

    open(tunnel) {
        if (this.tunnel) {
            this.tunnel.off('data', this.onTunnelData);
            this.tunnel.off('online', this.onTunnelOnline);
            this.tunnel.off('offline', this.onTunnelOffline);
            this.tunnel.off('close', this.onTunnelClose);
        }
        if (!this._checker)
            this._checker = timeout.check(1000, now => this.checkTimeout(now));
        tunnel.on('data', this.onTunnelData);
        tunnel.on('online', this.onTunnelOnline);
        tunnel.on('offline', this.onTunnelOffline);
        tunnel.on('close', this.onTunnelClose);

        this.tunnel = tunnel;
    }

    /**
     * 读取数据
     * @param {object} slave
     * @param {number} code
     * @param {number} address
     * @param {number} length
     * @param {boolean} quick
     * @returns {Promise<Buffer>}
     */
    read(slave, code, address, length, quick) {
        if (!this.tunnel.online) throw new Error("通道离线");
        if (quick) code = 4;
        const buf = Buffer.allocUnsafe(9);
        buf.writeUInt8(8)
        //buf.writeUInt8(this.transactionId);
        buf.writeUInt8(code, 2); //功能码
        buf.writeUInt8(slave.relay, 3); //中继
        buf.writeUInt8(slave.pool, 4); //塘号
        buf.writeUInt8(slave.device, 5); //设备
        buf.writeUInt8(address, 6); //地址
        buf.writeUInt8(length, 7); //长度
        //buf.writeUInt8(sum, 8);

        return this._execute(buf, quick);
    }

    /**
     * 写寄存器
     * @param {object} slave
     * @param {number} code
     * @param {number} address
     * @param {number} value
     * @returns {Promise<>}
     */
    write(slave, code, address, value) {
        if (!this.tunnel.online) throw new Error("通道离线");
        let buf = Buffer.allocUnsafe(10)
        buf.writeUInt8(9)
        //buf.writeUInt8(this.transactionId);
        buf.writeUInt8(5, 2); //功能码
        buf.writeUInt8(slave.relay, 3); //中继
        buf.writeUInt8(slave.pool, 4); //塘号
        buf.writeUInt8(slave.device, 5); //设备
        buf.writeUInt8(address, 6); //地址
        buf.writeUInt16BE(value, 7); //数据
        //buf.writeUInt8(sum, 9);

        return this._execute(buf, true);
    }

    _execute(command, primary) {
        //编号
        this.transactionId++
        if (this.transactionId > this.options.transaction.max)
            this.transactionId = this.options.transaction.min;
        command.writeUInt8(this.transactionId, 1);
        //和检验
        let sum = calcSum(command, command.length - 1);
        command.writeUInt8(sum & 0xFF, command.length - 1);

        //异步返回
        const promise = new Promise((resolve, reject) => {
            if (this._handlers[this.transactionId])
                this.reject(this.transactionId, new Error("ID被重用"));
            this._handlers[this.transactionId] = {id: this.transactionId, command, resolve, reject};
        });

        if (this.options.concurrency.enable) {
            if (this._doing < this.options.concurrency.max || primary) {
                this._doing++;
                this.tunnel.write(command);
                this._handlers[this.transactionId].stamp = Date.now();
            } else {
                this._queue.push(command);
            }
        } else {
            if (this._doing === 0) {
                this._doing++;
                this.tunnel.write(command);
                this._handlers[this.transactionId].stamp = Date.now();
            } else {
                if (primary) this._queue.splice(0, 0, command);
                else this._queue.push(command);
            }
        }

        return promise;
    }

    _next() {
        if (!this._queue.length) return;
        const cmd = this._queue.splice(0, 1)[0];
        this._doing++;
        this.tunnel.write(cmd);

        const id = cmd.readUInt8(1);
        this._handlers[id].stamp = Date.now();

        //console.log("next", cmd)
    }

    resolve(id, data) {
        const handler = this._handlers[id];
        delete this._handlers[id];
        if (handler)
            handler.resolve(data);
        if (this._doing > 0)
            this._doing--;
        process.nextTick(() => this._next());
    }

    reject(id, err) {
        const handler = this._handlers[id];
        delete this._handlers[id];
        if (handler)
            handler.reject(err);
        if (this._doing > 0)
            this._doing--;
        process.nextTick(() => this._next());
    }

    checkTimeout(now) {
        for (let i in this._handlers) {
            if (this._handlers.hasOwnProperty(i)) {
                const h = this._handlers[i];
                if (h.stamp && (h.stamp + this.options.timeout < now)) {
                    this.reject(h.id, new Error("超时了 " + h.command.toString("hex")));
                }
            }
        }
    }

    _handle(data) {
        if (data.length < 4) {
            //console.log("长度不能少于10字节")
            return;
        }

        //历史数据
        // if (this.remain) {
        //     data = Buffer.concat([this.remain, data]);
        //     this.remain = undefined;
        // }

        let len = 0;
        //定位正确的包
        while (data.length) {
            len = data[0];
            //过长，丢弃
            if (len > 100) {
                data = data.slice(1);
                continue;
            }
            //过短，也丢弃
            if (len < 4) {
                return;
            }

            //数据还不够 (考虑历史数据，可能会导致有效数据无法及时处理)
            // if (len >= data.length) {
            //     this.remain = data;
            //     return;
            // }

            if (len >= data.length) {
                data = data.slice(1);
                continue;
            }

            //校验
            let sum = calcSum(data, len) & 0xFF;
            let sum2 = data.readUInt8(len);
            if (sum !== sum2) {
                data = data.slice(1);
                continue;
            }

            break;
        }

        //如果找不到有效数据，就退出
        if (!data.length)
            return;

        let id = data.readUInt8(1);
        let fc = data.readUInt8(2);
        if ((fc & 0xF0) > 0) {
            let status = (fc & 0xF0) >> 4;
            let err;
            switch (status) {
                case 1:
                    err = "找不到中继";
                    break;
                case 2:
                    err = "找不到塘";
                    break;
                case 3:
                    err = "找不到设备";
                    break;
                case 4:
                    err = "地址错误";
                    break;
                case 5:
                    err = "长度错误";
                    break;
                case 6:
                    err = "校验错误";
                    break;
                case 7:
                    err = "功能码错误";
                    break;
                case 13:
                    err = "忙";
                    break;
                case 14:
                    err = "不忙";
                    break;
                case 15:
                    err = "解析错误";
                    break;
                default:
                    err = "未知错误";
                    break;
            }
            this.reject(id, new Error(err));
        } else {
            switch (fc) {
                case 1:  //心跳
                    break;
                case 3: //数据采集
                case 4: //读取
                {
                    let relay = data.readUInt8(3);
                    let pool = data.readUInt8(4);
                    let device = data.readUInt8(5);
                    let address = data.readUInt8(6);
                    this.resolve(id, data.slice(7, len));
                    break;
                }
                case 5: //写指令
                {
                    let relay = data.readUInt8(3);
                    let pool = data.readUInt8(4);
                    let device = data.readUInt8(5);
                    let address = data.readUInt8(6);
                    this.resolve(id, {relay, pool, device, address});
                    break;
                }
                default:
                    if (id > 0)
                        this.reject(id, new Error("不支持的FC:" + fc));
                    break;
            }
        }

        //粘包的情况
        if (data.length > len + 1) {
            data = data.slice(len + 1);
            //继续解析
            this._handle(data);
        }

    }
}