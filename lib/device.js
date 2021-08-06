const EventEmitter = require("events");
const collector = require("./collector");
const vm = require("vm");
const _ = require("lodash");

const influx = require_plugin("influxdb")

class Device extends EventEmitter {
    /**
     * @type Adapter
     */
    adapter;

    options = {};

    collectors = [];
    validators = [];
    scripts = [];

    variables = {};

    values = {
        a: {
            value: 1,
            last: Date.now(),
            watchers: []
        }
    };


    constructor(adapter, options) {
        super();
        this.adapter = adapter;
        Object.assign(this.options, options);

        this.options.variables.forEach(v => {
            let val = eval(v.default); //计算默认值
            Object.defineProperty(this.variables, v.name, {
                get() {
                    return val;
                },
                set(value) {
                    val = value;
                    //需要根据类型，转成Buffer
                    switch (v.type) {
                        case 'dword':
                            value = Buffer.allocUnsafe(4);
                            value.writeUInt32BE(val);
                            break;
                        case 'float':
                            value = Buffer.allocUnsafe(4);
                            value.writeFloatBE(val);
                            break;
                        case 'double':
                            value = Buffer.allocUnsafe(8);
                            value.writeDoubleBE(val);
                            break;
                        case 'uint32':
                            value = Buffer.allocUnsafe(4);
                            value.writeUInt32BE(val);
                            break;
                        case 'uint64':
                            value = Buffer.allocUnsafe(8);
                            value.writeBigUInt64BE(val);
                            break;
                        case 'int8':
                            value = Buffer.alloc(2);
                            value.writeInt8(val);
                            break;
                        case 'int16':
                            value = Buffer.allocUnsafe(2);
                            value.writeInt16BE(val);
                            break;
                        case 'int32':
                            value = Buffer.allocUnsafe(4);
                            value.writeInt32BE(val);
                            break;
                        case 'int64':
                            value = Buffer.allocUnsafe(8);
                            value.writeBigInt64BE(val);
                            break;
                        case 'le-float':
                            value = Buffer.allocUnsafe(4);
                            value.writeFloatLE(val);
                            break;
                        case 'le-double':
                            value = Buffer.allocUnsafe(8);
                            value.writeDoubleLE(val);
                            break;
                        case 'le-uint16':
                            value = Buffer.allocUnsafe(4);
                            value.writeUInt32LE(val);
                            break;
                        case 'le-uint32':
                            value = Buffer.allocUnsafe(4);
                            value.writeUInt32LE(val);
                            break;
                        case 'le-uint64':
                            value = Buffer.allocUnsafe(8);
                            value.writeBigUInt64LE(val);
                            break;
                        case 'le-int16':
                            value = Buffer.allocUnsafe(2);
                            value.writeInt16LE(val);
                            break;
                        case 'le-int32':
                            value = Buffer.allocUnsafe(4);
                            value.writeInt32LE(val);
                            break;
                        case 'le-int64':
                            value = Buffer.allocUnsafe(8);
                            value.writeBigInt64LE(val);
                            break;
                    }
                    adapter.write(options.slave, v.code, v.address, value)
                        .then(console.log).catch(console.error);
                }
            })
            //this.variables[v.name]
        })

        this.options.collectors.forEach(c => {
            const col = collector.create(adapter, options.slave, c);
            col.on("data", data => {

                const modify = [];
                const stores = [];
                const values = {};
                this.options.variables.forEach(v => {
                    if (v.code !== c.code) return;
                    if (v.address > c.length) return;
                    let value;
                    switch (v.type) {
                        case 'boolean':
                            value = data.readUInt8(v.address - c.address)
                            break;
                        case 'word':
                            value = data.readUInt16BE((v.address - c.address) * 2);
                            break;
                        case 'dword':
                            value = data.readUInt32BE((v.address - c.address) * 2);
                            break;
                        case 'float':
                            value = data.readFloatBE((v.address - c.address) * 2);
                            break;
                        case 'double':
                            value = data.readDoubleBE((v.address - c.address) * 2);
                            break;
                        case 'int8':
                            value = data.readInt8((v.address - c.address) * 2);
                            break;
                        case 'int16':
                            value = data.readInt16BE((v.address - c.address) * 2);
                            break;
                        case 'int32':
                            value = data.readInt32BE((v.address - c.address) * 2);
                            break;
                        case 'int64':
                            value = data.readBigInt64BE((v.address - c.address) * 2);
                            break;
                        case 'uint8':
                            value = data.readUInt8((v.address - c.address) * 2);
                            break;
                        case 'uint16':
                            value = data.readUInt16BE((v.address - c.address) * 2);
                            break;
                        case 'uint32':
                            value = data.readUInt32BE((v.address - c.address) * 2);
                            break;
                        case 'uint64':
                            value = data.readBigUInt64BE((v.address - c.address) * 2);
                            break;
                        case 'le-float':
                            value = data.readFloatLE((v.address - c.address) * 2);
                            break;
                        case 'le-double':
                            value = data.readDoubleLE((v.address - c.address) * 2);
                            break;
                        case 'le-int16':
                            value = data.readInt16LE((v.address - c.address) * 2);
                            break;
                        case 'le-int32':
                            value = data.readInt32LE((v.address - c.address) * 2);
                            break;
                        case 'le-int64':
                            value = data.readBigInt64LE((v.address - c.address) * 2);
                            break;
                        case 'le-uint16':
                            value = data.readUInt16LE((v.address - c.address) * 2);
                            break;
                        case 'le-uint32':
                            value = data.readUInt32LE((v.address - c.address) * 2);
                            break;
                        case 'le-uint64':
                            value = data.readBigUInt64LE((v.address - c.address) * 2);
                            break;
                    }
                    if (v.store) {
                        stores.push({
                            name: v.name,
                            type: v.type,
                            value: value
                        })
                    }
                    //构建合集
                    values[v.name] = value;
                    //记录变化
                    if (this.variables[v.name] !== value)
                        modify.push(v.name)
                })

                
                //统一赋值
                Object.assign(this.variables, values);

                //TODO 存入数据库
                if (stores.length) {
                    influx.write(options.type || options.name, [{name: 'id', value: options.value}], values)
                }

                //合法检查
                this.options.validators.forEach(v => {
                    const ret = vm.runInNewContext(v.expression, this.variables);
                    if (!ret) {
                        //不合法逻辑

                        //延时，发生时间，当前时间
                        //

                        //已经上报，则不再上报

                    } else {
                        //去掉发生时间，重置延时
                        //去掉已经上报标识

                    }
                })

                //向project汇报
                this.emit("values", values);

                //监听变化
                this.options.scripts.forEach(s => {
                    if (s.watches && s.watches.length && !_.intersection(s.watches, modify).length) return;
                    vm.runInNewContext(s.script, this.variables);
                    //TODO new vm.Script(s.script) 可以先编译代码提速
                })

            });
            this.collectors.push(col);
        })


    }

    execute(cmd, param) {
        const command = this.options.commands[cmd];
        vm.runInNewContext(command.script, this.variables);
        //TODO new vm.Script(command.script) 可以先编译代码提速
    }

}

exports.create = function (adapter, options) {

}