const timeout = require('./interval');

const crc16 = require("./crc16");

module.exports = class RTU {
    tunnel; //socket serial
    options = {
        timeout: 2000,
    };



    _handler;    //Promise
    
    _stamp = 0;
    
    _queue = [];

    _doing = false;
    
    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);

        const cancelable = timeout.check(1000, now => this.checkTimeout(now));

        tunnel.on('data', data => {
            this._handle(data)
        })

        tunnel.on('close', ()=>{
            cancelable.cancel()
        })
    }

    read(cmd) {
        const buf = Buffer.allocUnsafe(8);
        buf.writeUInt8(cmd.slave, 0);
        buf.writeUInt8(cmd.code, 1);
        buf.writeUInt16BE(cmd.address, 2);
        buf.writeUInt16BE(cmd.length, 4);
        buf.writeUInt16LE(crc16(buf.slice(0, -2)), 6);

        return this._execute(buf, cmd.primary);
    }

    write(cmd) {
        const buf = Buffer.allocUnsafe(6 + cmd.data.length);
        buf.writeUInt8(cmd.slave, 0);
        buf.writeUInt8(cmd.code, 1);
        buf.writeUInt16BE(cmd.address, 2);
        cmd.data.copy(buf, 4);
        buf.writeUInt16LE(crc16(buf.slice(0, -2)), buf.length - 2);

        return this._execute(buf, true);
    }


    _execute(command, primary) {
        let handler;
        //异步返回
        const promise = new Promise((resolve, reject) => {
            handler = {resolve, reject};
        });

        if (!this._doing) {
            this._doing = true;
            this.tunnel.write(command);
            this._handler = handler;
            this._stamp = Date.now();
        } else {
            //如果已经有指令在执行，则入队列等待
            const cmd = {
                command,
                handler
            };
            if (primary) this._queue.splice(0, 0, cmd);
            else this._queue.push(cmd);
        }

        return promise;
    }

    _next() {
        if (!this._queue.length) return;
        const cmd = this._queue.splice(0, 1)[0];
        this._doing = true;
        this.tunnel.write(cmd.command);
        this._stamp = Date.now();
        this._handler = cmd.handler;
        //console.log("next", cmd)
    }

    resolve(data) {
        if (this._handler) {
            this._handler.resolve(data);
            this._handler = undefined;
        }

        this._next();
    }

    reject(err) {
        if (this._handler) {
            this._handler.reject(err);
            this._handler = undefined;
        }

        this._next();
    }

    checkTimeout(now) {
        if (this._stamp + this.options.timeout < now) {
            this._doing = false;
            this.reject(new Error("超时了"));
        }
    }

    _handle(data) {
        this._doing = false;

        if (data.length < 5) {
            this.reject(new Error("长度不能少于5字节"));
            return;
        }

        const crc = data.readUInt16LE(data.length - 2);
        if (crc !== crc16(data.slice(0, -2))) {
            this.reject(new Error("检验错误"));
            return;
        }

        let slave = data.readUInt8(0);
        let fc = data.readUInt8(1);

        if ((fc & 0x80) > 0) {
            this.reject(new Error("执行错误 FC:" + fc + ' CODE:' + data.readUInt8(2)));
            return;
        }

        switch (fc) {
            case 1: //ReadCoils
            case 2: //ReadDiscreteInputs,
            {
                const count = data.readUInt8(2); // data[2];
                if (data.length - 5 !== count) {
                    this.reject(new Error("数据长度不对 COUNT:" + count + ", LENGTH:" + (data.length - 5)))
                    return;
                }

                //boolean数组展开
                let results = [];
                for (let i = 0; i < count; i++) {
                    let reg = data[i + 3];
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
                const count = data.readUInt8(2); //data[2];
                if (data.length - 5 !== count) {
                    this.reject(new Error("数据长度不对 COUNT:" + count + ", LENGTH:" + (data.length - 5)))
                    return;
                }

                let results = [];
                for (let i = 0; i < count; i += 2)
                    results.push(data.readUInt16BE(i + 3));
                this.resolve(results)
                break;
            }
            case 5: //WriteCoil
            case 6: //WriteHoldingRegister
            {
                const address = data.readUInt16BE(2);
                const value = data.readUInt16BE(4);
                this.resolve({address, value});
                break;
            }
            case 15: //WriteCoils
            case 16: //WriteHoldingRegisters
            {
                const address = data.readUInt16BE(2);
                const length = data.readUInt16BE(4);
                this.resolve({address, length});
                break;
            }
            default:
                break;
        }
    }
}
