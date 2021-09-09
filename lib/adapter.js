module.exports = class Adapter {
    protocol;
    slave;

    map = [];
    indexedMap = {};

    values = {};

    /**
     * @param {object} protocol
     * @param {object|number} slave
     * @param {array} map
     */
    constructor(protocol, slave, map) {
        this.protocol = protocol;
        this.slave = slave;

        this.map = map;
        //创建索引
        map.forEach(point => {
            this.indexedMap[point.name] = point;
            switch (point.type) {
                case 'dword':
                case 'int32':
                case 'uint32':
                    point.size = 2;
                    break;
                case 'double':
                case 'int64':
                case 'uint64':
                    point.size = 4;
                    break;
                default:
                    point.size = 1;
                    break;
            }
        });
    }

    /**
     * 设置
     * @param {string} key
     * @param {any} value
     * @returns {Promise<object>}
     */
    set(key, value) {
        const point = this.indexedMap[key]
        if (!point) throw new Error("未知数据点：" + key);
        const data = this._build(point, value)
        log.trace({key, value}, 'adapter set')
        return new Promise(((resolve, reject) => {
            this.protocol.write(this.slave, point.code, point.address, data).then(data => {
                log.trace({key, value}, 'adapter set ok')
                resolve(data);
            }).catch(reject);
        }))
    }

    /**
     * 获取
     * @param {string} key
     * @returns {Promise<any>}
     */
    get(key) {
        const point = this.indexedMap[key]
        if (!point) throw new Error("未知数据点：" + key);
        log.trace({key}, 'adapter get')
        return new Promise(((resolve, reject) => {
            this.protocol.read(this.slave, point.code, point.address, point.size, true).then(data => {
                const value = this._parse(point, data, 0)
                log.trace({key, value}, 'adapter get')
                resolve(value);
            }).catch(reject);
        }))
    }

    /**
     * 读取
     * @param {number} code
     * @param {number} address
     * @param {number} size
     * @param {boolean?} quick
     * @returns {Promise<object>}
     */
    read(code, address, size, quick) {
        log.trace({code, address, size}, 'adapter read')
        return new Promise(((resolve, reject) => {
            this.protocol.read(this.slave, code, address, size, quick).then(data => {
                const values = {}
                this.map.forEach(pt => {
                    if (pt.code !== code) return;
                    if (pt.address < address) return;
                    if (pt.address > address + size) return;
                    values[pt.name] = this._parse(pt, data, (pt.address - address) * 2)
                });
                log.trace(values, 'adapter read ok')
                resolve(values);
            }).catch(reject);
        }))
    }

    /**
     * 获得多个值
     * @param {string} key
     * @param {number} size
     * @param {boolean?} quick
     * @returns {Promise<object>}
     */
    getMany(key, size, quick) {
        const point = this.indexedMap[key]
        if (!point) throw new Error("未知数据点：" + key);
        return this.read(point.code, point.address, size, quick);
    }

    _build(point, val) {
        let value = val;
        switch (point.type) {
            case 'word':
                if (point.le) {
                    const buf = Buffer.allocUnsafe(2);
                    buf.writeUInt16LE(val);
                    value = buf.readUInt16BE();
                }
                break;
            case 'dword':
                value = Buffer.allocUnsafe(4);
                point.le ? value.writeUInt32LE(val) : value.writeUInt32BE(val);
                break;
            case 'float':
                value = Buffer.allocUnsafe(4);
                point.le ? value.writeFloatLE(val) : value.writeFloatBE(val);
                break;
            case 'double':
                value = Buffer.allocUnsafe(8);
                point.le ? value.writeDoubleLE(val) : value.writeDoubleBE(val);
                break;
            case 'uint8':
                if (point.le)
                    value *= 256;
                break;
            case 'uint16':
                if (point.le) {
                    const buf = Buffer.allocUnsafe(2);
                    buf.writeUInt16LE(val);
                    value = buf.readUInt16BE();
                }
                break;
            case 'uint32':
                value = Buffer.allocUnsafe(4);
                point.le ? value.writeUInt32LE(val) : value.writeUInt32BE(val);
                break;
            case 'uint64':
                value = Buffer.allocUnsafe(8);
                point.le ? value.writeInt8(val) : value.writeInt8(val);
                break;
            case 'int8':
                if (point.le)
                    value *= 256;
                break;
            case 'int16':
                if (point.le) {
                    const buf = Buffer.allocUnsafe(2);
                    buf.writeInt16LE(val);
                    value = buf.readInt16BE();
                }
                break;
            case 'int32':
                value = Buffer.allocUnsafe(4);
                point.le ? value.writeInt32LE(val) : value.writeInt32BE(val);
                break;
            case 'int64':
                value = Buffer.allocUnsafe(8);
                point.le ? value.writeBigInt64LE(val) : value.writeBigInt64BE(val);
                break;
        }
        return value;
    }

    _parse(point, data, offset = 0) {
        let value;
        switch (point.type) {
            case 'boolean':
                value = !!data.readUInt8(offset + 1);
                break;
            case 'word':
                value = point.le ? data.readUInt16LE(offset) : data.readUInt16BE(offset);
                break;
            case 'dword':
                value = point.le ? data.readUInt32LE(offset) : data.readUInt32BE(offset);
                break;
            case 'float':
                value = point.le ? data.readFloatLE(offset) : data.readFloatBE(offset);
                break;
            case 'double':
                value = point.le ? data.readDoubleLE(offset) : data.readDoubleBE(offset);
                break;
            case 'int8':
                value = data.readInt8(offset + 1);
                break;
            case 'int16':
                value = point.le ? data.readInt16LE(offset) : data.readInt16BE(offset);
                break;
            case 'int32':
                value = point.le ? data.readInt32LE(offset) : data.readInt32BE(offset);
                break;
            case 'int64':
                value = point.le ? data.readBigInt64LE(offset) : data.readBigInt64BE(offset);
                break;
            case 'uint8':
                value = data.readUInt8(offset + 1);
                break;
            case 'uint16':
                value = point.le ? data.readUInt16LE(offset) : data.readUInt16BE(offset);
                break;
            case 'uint32':
                value = point.le ? data.readUInt32LE(offset) : data.readUInt32BE(offset);
                break;
            case 'uint64':
                value = point.le ? data.readBigUInt64LE(offset) : data.readBigUInt64BE(offset);
                break;
            default:
                throw new Error("未知的数据类型" + point.type);
        }

        //精度计算
        if (point.type !== 'boolean' && point.precision > 0) {
            switch (point.precision) {
                case 1:
                    value *= 0.1;
                    break;
                case 2:
                    value *= 0.01;
                    break;
                case 3:
                    value *= 0.001;
                    break;
                case 4:
                    value *= 0.0001;
                    break;
                case 5:
                    value *= 0.00001;
                    break;
                case 6:
                    value *= 0.000001;
                    break;
            }
        }
        return value;
    }


}