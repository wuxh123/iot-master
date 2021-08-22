const EventEmitter = require("events");
const collector = require("./collector");
const _ = require("lodash");
const cron = require("./cron");
const project = require("./project");
const script = require("./script");
const {createEvent} = require("./event");
const Observer = require("./observer");

const influx = require_plugin("influxdb")
const mongo = require_plugin("mongodb")

class Device extends EventEmitter {
    /**
     * @type Tunnel
     */
    tunnel;

    agent;

    model = {};
    element = {};

    collectors = [];
    validators = [];
    jobs = [];
    scripts = [];

    commands = {};
    userJobs = {};

    //变量集，要使用defineProperty来定义变量
    variables = {$online: true};
    values = {};

    closed = true;
    error = '';

    constructor(tunnel, model) {
        super();
        this.open(tunnel, model);
    }

    onTunnelOnline = () => {
        this.variables.$online = true;
        this.variables.$onlineTime = Date.now();
        createEvent({device_id: this.model._id, event: '上线'})
        this.start();
    }

    onTunnelOffline = () => {
        this.variables.$online = false;
        this.variables.$offlineTime = Date.now();
        createEvent({device_id: this.model._id, event: '下线'})
        this.stop();
    }

    //箭头函数，避免this丢失
    onTunnelClose = () => {
        console.log('onTunnelClose', this.model._id, this.model.tunnel_id);

        //this.variables = {$online: false};//清空变量
        this.variables.$online = false; //不能清空

        //执行离线检查
        this.validate();

        //通知项目，设备离线
        this.emit('close');

        this.close();

        //记录下线
        createEvent({device_id: this.model._id, event: '下线'})
    }

    start() {
        if (this.closed) throw new Error('设备已经关闭');
        this.collectors.forEach(c => c.start());
    }

    stop() {
        this.collectors.forEach(c => c.cancel());
    }

    open(tunnel, model) {
        if (!this.closed)
            this.close();
        this.closed = false;

        this.tunnel = tunnel;
        if (model)
            this.model = model;

        //记录上线
        createEvent({device_id: this.model._id, event: '上线'})

        //通道关闭，则关闭设备（采集器，定时器等）
        tunnel.on('close', this.onTunnelClose);
        tunnel.on('online', this.onTunnelOnline);
        tunnel.on('offline', this.onTunnelOffline);

        this.variables.$online = true;
        this.variables.__proto__ = new Observer(); //直接继承

        //查找元件，以初始化
        mongo.db.collection("element").findOne({_id: model.element_id}).then(element => {
            this.element = element;

            //构建变量
            element.variables.forEach(v => this.variables.define(v.name, v.default && eval(v.default)));

            //数据点转为变量
            element.data_points.forEach(v => this.variables.define(v.name, v.default && eval(v.default)));

            //验证器
            element.validators.forEach(v => {
                if (!v.enable) return;
                this.validators.push({
                    script: script.compile(v.expression),
                    model: v,
                    start: 0,
                    delay: v.delay * 1000, //换算成s，避免后续冗余的乘法计算
                    resetInterval: v.resetInterval * 1000,
                    resetTimes: 0,
                });
            })

            //初始化命令
            element.commands.forEach(c => {
                this.commands[c.name] = {
                    script: script.compile(c.script),
                    model: c,
                };
            })

            //自动脚本
            element.scripts.forEach(s => {
                if (!s.enable) return;
                this.scripts.push({
                    script: script.compile(s.script),
                    model: s,
                });
            })

            //定时任务
            element.jobs.forEach(j => {
                if (!j.enable) return;

                const job = {script: script.compile(j.script), model: j};
                this.jobs.push(job)

                //启动定时
                job.handler = cron.schedule(j.crontab, () => {
                    try {
                        job.script.runInNewContext(this.variables);
                    } catch (err) {
                        this.error = err.message;
                        console.error(err)
                    }
                })
            })

            //创建采集器
            this.agent = tunnel.adapter.createAgent(model.slave, element.data_points);
            element.collectors.forEach(c => {
                if (!c.enable) return;
                this.addCollector(c);
            })


            //启动用户定时任务
            this.initUserJob();
        })
    }

    refresh() {
        this.collectors.forEach(c => c.execute())
    }

    close() {
        //不再监听关闭事件
        this.tunnel.off('close', this.onTunnelClose)
        this.tunnel.off('online', this.onTunnelOnline);
        this.tunnel.off('offline', this.onTunnelOffline);

        this.closed = true;
        //this.variables = {$online: false};//清空变量
        this.variables.$online = false;
        this.collectors.forEach(c => c.cancel());
        this.collectors = [];
        //this.validators = []; //检查器中有历史信息
        this.jobs.forEach(j => j.handler.cancel());
        this.jobs = [];
        this.scripts = [];

        this.clearUserJob();
    }

    addCollector(c) {
        const col = collector.create(this.agent, c);

        col.on("error", err => {
            //TODO 汇报给监控
            this.error = err.message
            //console.error(err)
            //this.emit("err", err);
        })

        col.on("data", values => {
            //找出要保存的
            const stores = this.element.data_points.filter(v => values.hasOwnProperty(v.name) && v.store).map(v => {
                return {
                    name: v.name,
                    //有倍率的值 从word转为float，避免数据库精度丢失
                    type: (v.ratio !== 0 && v.ratio !== 1 && v.type !== 'boolean') ? 'float' : v.type,
                    value: values[v.name]
                }
            });

            //找出变化的 key
            const modify = Object.keys(values).filter(k => this.variables[k] !== values[k]);

            //存入数据库
            if (stores.length) {
                //TODO 改用适配器，以兼容其他时序数据库
                const table = this.element.table || this.element._id;
                influx.write(table, [{name: 'id', value: this.model._id}], stores)
            }

            //统一赋值(直接赋原值，不会触发set，进而导致write)
            this.variables.set(values);

            //合法检查
            this.validate();

            //向project汇报
            this.emit("values", values);

            //监听变化
            if (modify.length) {
                this.scripts.forEach(s => {
                    if (s.model.watches && s.model.watches.length && !_.intersection(s.model.watches, modify).length)
                        return;
                    try {
                        s.script.runInNewContext(this.variables);
                    } catch (err) {
                        this.error = err.message;
                        console.error(err)
                    }
                })

                //向project汇报
                this.emit("modify", modify);
            }
        });

        //添加之后就执行一次采集
        col.execute();

        this.collectors.push(col);
    }

    initUserJob() {
        //启动用户定时任务
        mongo.db.collection("job").find({device_id: this.model._id, enable: true}).toArray().then(jobs => {
            jobs.forEach(j => this.addUserJob(j));
        }).catch(err => {
            this.error = err.message;
            console.error(err)
        });
    }

    addUserJob(j) {
        const job = {model: j};
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
                    params = job.script.runInNewContext(Object.assign({}, this.variables, this.variables.values()));
                }

                this.execute(j.command, params).then(() => {
                    createEvent({device_id: this.model._id, event: '定时任务：' + j.command})
                }).catch(err => {
                    createEvent({device_id: this.model._id, event: '定时任务：' + j.command + ' 失败：' + err.message})
                });
            } catch (err) {
                this.error = err.message;
                console.error(err)
            }
        })
    }

    clearUserJob() {
        Object.keys(this.userJobs).forEach(k => {
            const job = this.userJobs[k];
            job.handler.cancel();
        })
        this.userJobs = {};
    }

    removeUserJob(_id) {
        if (this.userJobs[_id]) {
            this.userJobs[_id].handler.cancel();
            delete this.userJobs[_id];
        }
    }


    validate() {
        this.validators.forEach(v => {
            try {
                const ctx = Object.assign({}, this.variables, this.variables.values())
                const ret = v.script.runInNewContext(ctx);
                if (ret) {
                    v.start = 0;//去掉发生时间，重置延时
                    v.reported = false;//去掉已经上报标识
                    return;
                }
            } catch (err) {
                this.error = err.message;
                console.error(err)
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
            if (v.reported) {
                //TODO 重置逻辑
                if (v.model.resetInterval && v.model.resetTimes) {
                    //如果已经超出重置次数，则不再执行
                    if (v.resetTimes > v.model.resetTimes)
                        return;
                    //如果还没到重置时间，则不提醒
                    if (v.reportAt + v.resetInterval < now)
                        return;
                    //重置
                    v.start = now;
                    v.resetTimes++
                    return; //下次再执行，因为可能会有delay
                }
                return;
            }
            v.reported = true;
            v.reportAt = Date.now()

            //TODO content 也可以使用表达式（脚本），但是变量要先拷贝，避免被覆盖
            this.emit('alarm', {
                name: v.model.name,
                content: v.model.content,
                level: v.model.level,
                device_id: this.model._id, //TODO device name
                device_name: this.model.name || this.element.name, //添加设备名
            });
        })
    }

    async writeValues(values) {
        for (let k in values) {
            if (!values.hasOwnProperty(k)) return;
            await this.agent.set(k, values[k])
        }
    }

    async execute(cmd, param) {
        if (this.closed)
            throw new Error("设备离线")

        const command = this.commands[cmd];
        if (!command)
            throw new Error("不支持的命令")

        //创建子对象，避免污染variables
        const ctx = {};
        Array.isArray(param) && param.forEach((v, i) => ctx['$' + (i + 1)] = v);
        ctx.__proto__ = this.variables;
        command.script.runInNewContext(ctx);
        let changes = this.variables.changes();
        //console.log('execute changes', changes)

        //将修改值存入
        await this.writeValues(changes);
        //设备离线，可以缓存一段时间？

        //重新采集一次数据
        this.refresh();
    }
}

const devices = {};

/**
 * 打开设备（或恢复）
 * @param {Tunnel} tunnel
 * @param {Object} model
 * @returns {Device}
 */
exports.open = function (tunnel, model) {
    //离线恢复 逻辑
    let device = devices[model._id]
    if (device) {
        //device.open(tunnel, model);
        //device.start();
    } else {
        device = new Device(tunnel, model);
        //建立索引
        devices[model._id] = device;
    }

    //将设备放入项目中
    mongo.db.collection("project").find({
        enable: true,
        devices: {$elemMatch: {device_id: model._id}}
    }).toArray().then(projects => {
        projects.forEach(p => {
            const prj = project.get(p._id);
            if (prj) {
                prj.setDevice(model._id, device);
            } else {
                //TODO 恢复项目

            }
        })
    })

    return device;
}

/**
 * 获得设备
 * @param {string} id
 * @returns {Device}
 */
exports.get = function (id) {
    return devices[id];
}

/**
 * 删除设备
 * @param {string} id
 */
exports.remove = function (id) {
    const device = devices[id];
    if (device) {
        device.close();
        delete devices[id];  //devices[id]=null
    }
}
