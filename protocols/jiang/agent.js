module.exports = class Agent {
    adapter;
    slave;

    map = [];
    indexedMap = {};

    values = {};

    constructor(adapter, slave, map) {
        this.adapter = adapter;
        this.slave = slave;

        this.map = map;
        //创建索引
        map.forEach(point => {
            this.indexedMap[point.name] = point;
        });
    }

    set(key, value) {
        const point = this.indexedMap[key]
        if (!point) throw new Error("未知数据点：" + key);
        return new Promise(((resolve, reject) => {
            this.adapter.write(this.slave, point.code, point.address, value).then(data => {
                resolve(data);
            }).catch(reject);
        }))
    }

    get(key) {
        const point = this.indexedMap[key]
        if (!point) throw new Error("未知数据点：" + key);
        return new Promise(((resolve, reject) => {
            this.adapter.read(this.slave, point.code, point.address, 1).then(data => {
                resolve(this._parse(point, data, 0));
            }).catch(reject);
        }))
    }

    read(code, address, size) {
        return new Promise(((resolve, reject) => {
            this.adapter.read(this.slave, code, address, size).then(data => {
                const values = {}
                this.map.forEach(pt => {
                    //if (pt.code !== code) return;
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

    _parse(point, data, offset = 0) {
        let value;
        switch (point.type) {
            case 'word':
            case 'uint16':
                value = data.readUInt16BE(offset);
                break;
            case 'int16':
                value = data.readInt16BE(offset);
                break;
            default:
                throw new Error("未知的数据类型" + point.type);
        }

        //精度计算
        if (point.precision > 0) {
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