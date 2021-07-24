const timeout = require('./interval');
const helper = require('./helper');

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

        tunnel.on('close', () => {
            cancelable.cancel()
        })
    }

    /**
     * 读取数据
     * @param {number} slave
     * @param {string} address
     * @param {number} length
     * @returns {Promise<Buffer>}
     */
    read(slave, address, length) {
        let {code, address} = helper.parseReadAddress(address);
        const buf = Buffer.allocUnsafe(8);
        buf.writeUInt8(slave, 0); //从站号
        buf.writeUInt8(code, 1); //功能码
        buf.writeUInt16BE(address, 2); //地址
        buf.writeUInt16BE(length, 4); //长度
        buf.writeUInt16LE(helper.crc16(buf.slice(0, -2)), 6); //校验

        return this._execute(buf, false);
    }

    /**
     * 写单个线圈或寄存器
     * @param {number} slave
     * @param {string} address
     * @param {boolean|number} value
     * @returns {Promise<>}
     */
    write(slave, address, value) {
        let {code, address} = helper.parseWriteAddress(address);
        const buf = Buffer.allocUnsafe(8);
        buf.writeUInt8(slave, 0); //从站号
        buf.writeUInt8(code, 1); //功能码
        buf.writeUInt16BE(address, 2); //地址
        if (code === 5)
            buf.writeUInt16BE(value ? 0xFF00 : 0x0000, 4); //写线圈，0xFF00代表合，0x0000代表开
        else
            buf.writeUInt16BE(value, 4);
        buf.writeUInt16LE(helper.crc16(buf.slice(0, -2)), buf.length - 2); //检验位

        return this._execute(buf, true);
    }

    /**
     * 写入多个线圈或寄存器
     * @param {number} slave
     * @param {string} address
     * @param {boolean[]|Uint8Array|Uint16Array|Buffer} data
     * @returns {Promise<>}
     */
    writeMany(slave, address, data) {
        let {code, address} = helper.parseWriteAddress(address);
        code += 10; // 5=>15 6=>16

        let buffer;
        if (code === 15) {
            buffer = helper.compressBooleans(data)
        } else if (code === 16) {
            buffer = helper.arrayToBuffer(data)
        }

        const buf = Buffer.allocUnsafe(8 + buffer.length);
        buf.writeUInt8(slave, 0); //从站号
        buf.writeUInt8(code, 1); //功能码
        buf.writeUInt16BE(address, 2); //地址
        buf.writeUInt16BE(data.length, 4); //长度
        buffer.copy(buf, 6); //内容
        buf.writeUInt16LE(helper.crc16(buf.slice(0, -2)), buf.length - 2); //检验位

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
        if (crc !== helper.crc16(data.slice(0, -2))) {
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
                        results.push((reg & 1) === 1); // ? 1 : 0);
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

                // let results = [];
                // for (let i = 0; i < count; i += 2)
                //     results.push(data.readUInt16BE(i + 3));
                // 直接返回Buffer，方便外部解析非WORD类型，比如：浮点数，字符串
                this.resolve(data.slice(3, -2)); //count*2
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
