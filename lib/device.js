const EventEmitter = require("events");
const collector = require("./collector");
const _ = require("lodash");
const cron = require("./cron");
const project = require("./project");
const script = require("./script");

const influx = require_plugin("influxdb")
const mongo = require_plugin("mongodb")

class Device extends EventEmitter {
    /**
     * @type Tunnel
     */
    tunnel;

    model = {};
    element = {};

    collectors = [];
    validators = [];
    jobs = [];
    scripts = [];

    commands = {};
    userJobs = {};

    //变量集，要使用defineProperty来定义变量
    variables = {};
    values = {};

    closed = true;
    error = '';

    constructor(tunnel, model) {
        super();
        this.open(tunnel, model);
    }

    open(tunnel, model) {
        this.tunnel = tunnel;
        if (model)
            this.model = model;

        if (!this.closed)
            this.close();
        this.closed = false;

        //通道关闭，则关闭设备（采集器，定时器等） TODO 设置关闭状态 或 从索引中移除
        tunnel.on('close', () => {
            this.close();
        })

        //查找元件，以初始化
        mongo.db.collection("element").findOne({_id: model.element_id}).then(element => {
            this.element = element;

            //构建变量
            element.variables.forEach(v => {
                this.addVariable(v);
            })

            element.validators.forEach(v => {
                this.addValidator(v);
            })

            element.collectors.forEach(c => {
                if (c.enable) this.addCollector(c);
            })

            element.commands.forEach(c => {
                this.addCommand(c);
            })

            element.scripts.forEach(s => {
                if (s.enable) this.addScript(s);
            })

            element.jobs.forEach(j => {
                if (j.enable) this.addJob(j);
            })

            //启动用户定时任务
            mongo.db.collection("job").find({device_id: this.model._id, enable: true}).toArray().then(jobs => {
                jobs.forEach(j => this.addUserJob(j));
            }).catch(err => {
                this.error = err.message
            });
        })
    }

    refresh(){
        this.collectors.forEach(c=>c.execute())
    }

    close() {
        this.closed = true;
        this.collectors.forEach(c => c.cancel());
        this.collectors = [];
        this.jobs.forEach(j => j.handler.cancel());
        this.jobs = [];
        for (let id in this.userJobs) {
            if (this.userJobs.hasOwnProperty(id))
                this.userJobs[id].handler.cancel()
        }
        this.userJobs = {};
    }

    addVariable(v) {
        const that = this;
        that.values[v.name] = eval(v.default); //计算默认值
        Object.defineProperty(this.variables, v.name, {
            enumerable: true,
            get() {
                return that.values[v.name];
            },
            set(value) {
                that.values[v.name] = value;
                //需要类型，转换
                const data = that.tunnel.adapter.buildData(v.type, value);
                that.tunnel.adapter.write(that.model.slave, v.code, v.address, data)
                    .then(console.log)
                    .catch(err => {
                        this.error = err.message
                    });
                //TODO 报出设备错误，写日志
            }
        })
    }

    addCollector(c) {
        const col = collector.create(this.tunnel.adapter, this.model.slave, c);

        col.on("error", error => {
            //TODO 汇报给监控
            this.error = error.message
            //this.emit("error", error);
        })

        col.on("data", data => {

            //解析数据
            const values = this.tunnel.adapter.parseData(this.element.variables, data, c.code, c.address, c.length);

            //找出要保存的
            const stores = this.element.variables.filter(v => values.hasOwnProperty(v.name) && v.store).map(v => {
                return {
                    name: v.name,
                    //有倍率的值 从word转为float，避免数据库精度丢失
                    type: (v.ratio !== 0 && v.ratio !== 1 && v.type !== 'boolean') ? 'float' : v.type,
                    value: values[v.name]
                }
            });

            //找出变化的 key
            const modify = Object.keys(values).filter(k => this.variables[k] !== values[k]);

            //统一赋值(直接赋原值，不会触发set，进而导致write)
            Object.assign(this.values, values);

            //存入数据库
            if (stores.length) {
                //TODO 改用适配器，以兼容其他时序数据库
                const table = this.element.table || this.element._id;
                influx.write(table, [{name: 'id', value: this.model._id}], stores)
            }

            //合法检查
            this.validate();

            //向project汇报
            this.emit("values", values);

            //监听变化
            if (modify.length) {
                this.scripts.forEach(s => {
                    if (s.model.watches && s.model.watches.length && !_.intersection(s.model.watches, modify).length)
                        return;
                    s.script.runInNewContext(this.variables);
                })

                //向project汇报
                this.emit("modify", modify);
            }

        });

        //添加之后就执行一次采集
        col.execute();

        this.collectors.push(col);
    }

    addValidator(v) {
        this.validators.push({
            script: script.compile(v.expression),
            model: v,
            start: 0,
            delay: v.delay * 1000 //换算成s，避免后续冗余的乘法计算
        })
    }

    addCommand(c) {
        this.commands[c.name] = {
            script: script.compile(c.script),
            model: c,
        };
    }

    addJob(j) {
        const job = {
            script: script.compile(j.script),
            model: j,
        };
        this.jobs.push(job)

        //启动定时
        job.handler = cron.schedule(j.crontab, () => {
            //job.script.runInNewContext(this.variables);
            try {
                job.script.runInNewContext(this.variables);
            } catch (err) {
                this.error = err.message;
            }
        })
    }

    addUserJob(j) {
        const job = {
            model: j,
        };
        //this.jobs.push(job);
        if (this.userJobs[j._id])
            this.userJobs[j._id].handler.cancel();
        this.userJobs[j._id] = job;

        //预编译参数数组
        if (j.params)
            job.script = script.compile('[' + j.params + ']')

        //启动定时
        job.handler = cron.scheduleTime(j.time, () => {
            try {
                let params = [];
                if (j.params) {
                    //params = job.script.runInNewContext(Object.assign({}, this.variables)); //隔离，避免操作（项目应该深度clone）
                    params = job.script.runInNewContext(_.cloneDeep(this.variables));
                }
                this.execute(j.command, params)
            } catch (err) {
                this.error = err.message;
            }
        })
    }

    removeUserJob(_id) {
        if (this.userJobs[_id]) {
            this.userJobs[_id].handler.cancel();
            delete this.userJobs[_id];
        }
    }

    addScript(c) {
        this.scripts.push({
            script: script.compile(c.script),
            model: c,
        });
    }

    validate() {
        this.validators.forEach(v => {
            //const ret = vm.runInNewContext(v.expression, this.variables);
            try {
                const ret = v.script.runInNewContext(_.cloneDeep(this.variables));
                if (ret) {
                    v.start = 0;//去掉发生时间，重置延时
                    v.reported = false;//去掉已经上报标识
                    return;
                }
            } catch (err) {
                this.error = err.message;
                return;
            }
            //以下是不合法处理逻辑

            //延时处理，发生时间，当前时间
            const now = Date.now();
            if (v.delay) {
                if (!v.start) {
                    v.start = now;
                    return;
                }
                if (v.start + v.delay < now)
                    return;
            }

            //已经上报，则不再上报
            if (v.reported)
                return;

            //TODO content 也可以使用表达式（脚本），但是变量要先拷贝，避免被覆盖
            this.emit('alarm', v);
        })
    }


    execute(cmd, param) {
        //TODO 记录执行了命令

        const command = this.commands[cmd];
        if (!command) {
            throw new Error("不支持的命令")
        }

        //以下代码会污染variables，具体再定
        Array.isArray(param) && param.forEach((v, i) => this.variables['$' + (i + 1)] = v);
        try {
            command.script.runInNewContext(this.variables);
            //已验证 Object.assign({}, this.variables, params)作为context，set失效
            //每次都构造context，可能太低效了
        } catch (err) {
            this.error = err.message;
        }

        //设备离线，可以缓存一段时间？

        //重新采集一次数据
        this.refresh();
    }
}

const devices = {};

exports.create = function (tunnel, model) {
    //TODO 离线恢复 逻辑

    const device = new Device(tunnel, model);

    //建立索引
    devices[model._id] = device;

    //启动相关项目
    //放到接收器中
    mongo.db.collection("project").find({
        enable: true,
        devices: {$elemMatch: {device_id: model._id}}
    }).toArray().then(projects => {
        projects.forEach(p => {
            const prj = project.get(p._id);
            if (prj) {
                prj.addDevice(model._id, device);
            } else {
                //TODO 恢复项目
            }
        })
    })

    return device;
}

exports.get = function (id) {
    return devices[id];
}

exports.remove = function (id) {
    const device = devices[id];
    if (device) {
        device.close();
        delete devices[id];  //devices[id]=null
    }
}
