module.exports = class Adapter {
    protocol;
    slave;

    map = [];
    indexedMap = {};

    values = {};

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

    set(key, value) {
        const point = this.indexedMap[key]
        if (!point) throw new Error("未知数据点：" + key);
        const data = this._build(point, value)
        //console.log('Agent set', key, value, new Date());
        return new Promise(((resolve, reject) => {
            this.protocol.write(this.slave, point.code, point.address, data).then(data => {
                //console.log('Agent set', key, value, 'ok',new Date());
                resolve(data);
            }).catch(reject);
        }))
    }

    get(key) {
        const point = this.indexedMap[key]
        if (!point) throw new Error("未知数据点：" + key);
        return new Promise(((resolve, reject) => {
            this.protocol.read(this.slave, point.code, point.address, point.size).then(data => {
                resolve(this._parse(point, data, 0));
            }).catch(reject);
        }))
    }

    read(code, address, size) {
        return new Promise(((resolve, reject) => {
            this.protocol.read(this.slave, code, address, size).then(data => {
                const values = {}
                this.map.forEach(pt => {
                    if (pt.code !== code) return;
                    if (pt.address < address) return;
                    if (pt.address > address + size) return;
                    values[pt.name] = this._parse(pt, data, (pt.address - address) * 2)
                });
                resolve(values);
            }).catch(reject);
        }))
    }

    getMany(key, size) {
        const point = this.indexedMap[key]
        if (!point) throw new Error("未知数据点：" + key);
        return this.read(point.code, point.address, size);
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