const timeout = require('../../lib/interval');
const Agent = require("./agent");

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
    }

    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
        this.open(tunnel)
    }

    open(tunnel) {
        if (this.tunnel) {
            this.tunnel.off('data', this.onTunnelData);
            this.tunnel.off('close', this.onTunnelClose);
        }
        if (!this._checker)
            this._checker = timeout.check(1000, now => this.checkTimeout(now));
        tunnel.on('data', this.onTunnelData);
        tunnel.on('close', this.onTunnelClose);
    }

    /**
     * 创建Agent
     * @param {number} slave
     * @param {Object[]} map
     * @returns {Agent}
     */
    createAgent(slave, map) {
        return new Agent(this, slave, map);
    }

    /**
     * 读取数据
     * @param {number} slave
     * @param {number} code
     * @param {number} address
     * @param {number} length
     * @returns {Promise<Buffer>}
     */
    read(slave, code, address, length) {
        const buf = Buffer.allocUnsafe(9);
        buf.writeUInt8(8)
        //buf.writeUInt8(this.transactionId);
        buf.writeUInt8(code, 2); //功能码
        buf.writeUInt8((slave >> 16) & 0xFF, 3); //中继
        buf.writeUInt8((slave >> 8) & 0xFF, 4); //塘号
        buf.writeUInt8(slave & 0xFF, 5); //设备
        buf.writeUInt8(address, 6); //地址
        buf.writeUInt8(length, 7); //长度
        //buf.writeUInt8(sum, 8);

        return this._execute(buf, false);
    }

    /**
     * 写寄存器
     * @param {number} slave
     * @param {number} code
     * @param {number} address
     * @param {number} value
     * @returns {Promise<>}
     */
    write(slave, code, address, value) {
        let buf = Buffer.allocUnsafe(10)
        buf.writeUInt8(9)
        //buf.writeUInt8(this.transactionId);
        buf.writeUInt8(5, 2); //功能码
        buf.writeUInt8((slave >> 16) & 0xFF, 3); //中继
        buf.writeUInt8((slave >> 8) & 0xFF, 4); //塘号
        buf.writeUInt8(slave & 0xFF, 5); //设备
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
        command.writeUInt8(sum, command.length - 1);

        //异步返回
        const promise = new Promise((resolve, reject) => {
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
        const handler = this._handlers[id]
        if (handler) {
            handler.resolve(data);
            delete this._handlers[id]
        }

        this._next();
    }

    reject(id, err) {
        const handler = this._handlers[id]
        if (handler) {
            handler.reject(err);
            delete this._handlers[id]
        }

        this._next();
    }

    checkTimeout(now) {
        for (let i in this._handlers) {
            if (this._handlers.hasOwnProperty(i)) {
                const h = this._handlers[i];
                if (h.stamp && (h.stamp + this.options.timeout < now)) {
                    if (this._doing > 0)
                        this._doing--;
                    this.reject(h.id, new Error("超时了 " + h.command.toString("hex")));
                }
            }
        }
    }

    _handle(data) {
        if (this._doing > 0)
            this._doing--;

        if (data.length < 4) {
            //console.log("长度不能少于10字节")
            return;
        }

        //历史数据
        if (this.remain) {
            data = Buffer.concat([this.remain, data]);
            this.remain = undefined;
        }

        //定位正确的包
        while (true) {
            let len = data[0];
            //过长，丢弃
            if (len > 50) {
                data = data.slice(1);
                continue;
            }

            //数据还不够
            if (len >= data.length) {
                this.remain = data;
                return;
            }

            //校验
            let sum = calcSum(data, len);
            let sum2 = data.readUInt8(len);
            if (sum !== sum2) {
                data = data.slice(1);
                continue;
            }

            break;
        }

        let id = data.readUInt8(1);
        let fc = data.readUInt8(2);
        if ((fc & 0xF0) > 0) {
            let status = (fc & 0xF0) >> 4;
            this.reject(id, new Error("执行错误 FC:" + fc + ' STATUS:' + status));
            return;
        }

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
                this.resolve(id, data.slice(7, data.length - 1));
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
                break;
        }
    }
}