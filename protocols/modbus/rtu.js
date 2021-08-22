const timeout = require('./interval');
const helper = require('./helper');
const Agent = require("./agent");

module.exports = class RTU {

    tunnel; //socket serial
    options = {
        timeout: 2000,
    };


    _handler;    //Promise

    _stamp = 0;

    _queue = [];

    _doing = false;


    _checker;

    onTunnelData= data=>{
        this._handle(data)
    }
    onTunnelClose= ()=>{
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
        const buf = Buffer.allocUnsafe(8);
        buf.writeUInt8(slave, 0); //从站号
        buf.writeUInt8(code, 1); //功能码
        buf.writeUInt16BE(address, 2); //地址
        buf.writeUInt16BE(length, 4); //长度
        buf.writeUInt16LE(helper.crc16(buf.slice(0, -2)), 6); //校验

        return this._execute(buf, false);
    }

    /**
     * 写线圈或寄存器
     * @param {number} slave
     * @param {number} code
     * @param {number} address
     * @param {boolean|number|boolean[]|Uint8Array|Uint16Array|Buffer} value
     * @returns {Promise<>}
     */
    write(slave, code, address, value) {
        const type = typeof value;

        let buf;
        if (type === 'boolean' || type === 'number') {
            code = helper.convertWriteCode(code);
            buf = Buffer.allocUnsafe(8);
            if (code === 5)
                buf.writeUInt16BE(value ? 0xFF00 : 0x0000, 4); //写线圈，0xFF00代表合，0x0000代表开
            else
                buf.writeUInt16BE(value, 4);
        } else {
            code = helper.convertWriteCode(code, true);
            let buffer;
            if (code === 15) {
                buffer = helper.booleanArrayToBuffer(value)
                buf = Buffer.allocUnsafe(8 + buffer.length);
                buf.writeUInt16BE(value.length, 4);
                buffer.copy(buf, 6); //内容
            } else if (code === 16) {
                buffer = helper.arrayToBuffer(value)
                buf = Buffer.allocUnsafe(8 + buffer.length);
                buf.writeUInt16BE((buffer.length - 1) / 2, 4);
                buffer.copy(buf, 6); //内容
            } else {
                //暂时不会发生
            }
        }

        //包头和尾
        buf.writeUInt8(slave, 0); //从站号
        buf.writeUInt8(code, 1); //功能码
        buf.writeUInt16BE(address, 2); //地址
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
        if (this._doing && this._stamp + this.options.timeout < now) {
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
                        results.push((reg & 1) === 1 ? 1 : 0);
                        reg = reg >> 1;
                    }
                }

                //转成双字节码，方便解析
                //results = Buffer.from(new Uint16Array(results))
                const buf = Buffer.allocUnsafe(count * 8 * 2);
                results.forEach((r, i) => buf.writeUInt16BE(r, i * 2));

                this.resolve(buf)
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
