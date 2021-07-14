const timeout = require('./interval');

module.exports = class TCP {

    tunnel; //socket serial
    options = {
        concurrency: {
            enable: true,
            max: 10
        },
        timeout: 5000,
        transaction: {
            min: 0x5a01,
            max: 0x5aff
        }
    };

    transactionId = 0x5a01;

    _handlers = {};

    _queue = [];

    _doing = 0;


    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);

        const cancelable = timeout.check(2000, now => this.checkTimeout(now));

        tunnel.on('data', data => {
            this._handle(data)
        })

        tunnel.on('close', () => {
            cancelable.cancel()
        })
    }

    read(cmd) {
        const buf = Buffer.allocUnsafe(12);
        //buf.writeUInt16BE(this.transactionId);
        buf.writeUInt16BE(0, 2);
        buf.writeUInt16BE(6, 4);
        buf.writeUInt8(cmd.slave, 6);
        buf.writeUInt8(cmd.code, 7);
        buf.writeUInt16BE(cmd.address, 8);
        buf.writeUInt16BE(cmd.length, 10);

        return this._execute(buf, cmd.primary);
    }

    write(cmd) {
        const buf = Buffer.allocUnsafe(10 + cmd.data.length);
        //buf.writeUInt16BE(this.transactionId);
        buf.writeUInt16BE(0, 2);
        buf.writeUInt16BE(6, 4);
        buf.writeUInt8(cmd.slave, 6);
        buf.writeUInt8(cmd.code, 7);
        buf.writeUInt16BE(cmd.address, 8);
        cmd.data.copy(buf, 10);

        return this._execute(buf, true);
    }

    _execute(command, primary) {
        this.transactionId++
        if (this.transactionId > this.options.transaction.max)
            this.transactionId = this.options.transaction.min;
        command.writeUInt16BE(this.transactionId);

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

        const id = cmd.readUInt16BE();
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

        if (data.length < 12) {
            // if (data.length >= 2) {
            //     let id = data.readUInt16BE(0);
            //     if (id < 512) {
            //         this.reject(id, new Error("长度不能少于12字节"));
            //     }
            // }
            console.log("长度不能少于12字节")
            return;
        }

        let id = data.readUInt16BE(0);
        let protocol = data.readUInt16BE(2);
        let length = data.readUInt16BE(4);
        if (data.length < length + 6) {
            this.reject(id, new Error("长度不够"));
            //TODO 如果长度不够，需要等待
            return;
        }

        let slave = data.readUInt8(6);
        let fc = data.readUInt8(7);

        if ((fc & 0x80) > 0) {
            this.reject(id, new Error("执行错误 FC:" + fc + ' CODE:' + data.readUInt8(8)));
            return;
        }

        switch (fc) {
            case 1: //ReadCoils
            case 2: //ReadDiscreteInputs,
            {
                const count = data.readUInt8(8);
                if (length - 3 !== count) {
                    this.reject(id, new Error("数据长度不对 COUNT:" + count + ", LENGTH:" + (length - 3)))
                    return;
                }

                //boolean数组展开
                let results = [];
                for (let i = 0; i < count; i++) {
                    let reg = data[i + 9];
                    for (let j = 0; j < 8; j++) {
                        results.push((reg & 1) === 1);
                        reg = reg >> 1;
                    }
                }
                this.resolve(results)
                break;
            }
            case 3://ReadHoldingRegisters,
            case 4://ReadInputRegisters,
                //case 23://ReadWriteMultipleRegisters:
            {
                const count = data.readUInt8(8); //data[2];
                if (length - 3 !== count) {
                    this.reject(id, new Error("数据长度不对 COUNT:" + count + ", LENGTH:" + (length - 3)))
                    return;
                }

                let results = [];
                for (let i = 0; i < count; i += 2)
                    results.push(data.readUInt16BE(i + 9));
                this.resolve(id, results)
                break;
            }
            case 5: //WriteCoil
            case 6: //WriteHoldingRegister
            {
                const address = data.readUInt16BE(8);
                const value = data.readUInt16BE(10);
                this.resolve(id, {address, value});
                break;
            }
            case 15: //WriteCoils
            case 16: //WriteHoldingRegisters
            {
                const address = data.readUInt16BE(8);
                const length = data.readUInt16BE(10);
                this.resolve(id, {address, length});
                break;
            }
            default:
                break;
        }
    }
}