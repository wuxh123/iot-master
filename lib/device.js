const EventEmitter = require("events");
const _ = require("lodash");
const cron = require("./cron");
const project = require("./project");
const script = require("./script");
const {createEvent} = require("./event");
const Observer = require("./observer");
const interval = require("./interval");
const Adapter = require("./adapter");
const Validator = require("./validator");
const Collector = require("./collector");
const Job = require("./job");
const UserJob = require("./user_job");

const influx = require_plugin("influxdb")
const mongo = require_plugin("mongodb")

class Device extends EventEmitter {
    /**
     * @type Tunnel
     */
    tunnel;

    adapter;

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

        this.on('error', err => {
            this.error = err.message;
            log.error({id: this.model._id, error: err.message}, '设备错误')
        })

        this.open(tunnel, model);
    }

    onTunnelOnline = () => {
        log.info({id: this.model._id}, '设备上线');
        this.variables.$online = true;
        this.variables.$onlineTime = Date.now();
        createEvent({device_id: this.model._id, event: '上线'})
        this.start();
    }

    onTunnelOffline = () => {
        log.info({id: this.model._id}, '设备下线');
        this.variables.$online = false;
        this.variables.$offlineTime = Date.now();
        createEvent({device_id: this.model._id, event: '下线'})
        this.stop();
    }

    //箭头函数，避免this丢失
    onTunnelClose = () => {
        //console.log('onTunnelClose', this.model._id, this.model.tunnel_id);
        log.info('onTunnelClose')

        //this.variables = {$online: false};//清空变量
        this.variables.$online = false; //不能清空

        //执行离线检查
        //this.validate();

        //通知项目，设备离线
        this.emit('close');

        this.close();
    }

    start() {
        log.info({id: this.model._id}, '启动设备');
        if (this.closed) throw new Error('设备已经关闭');
        this.collectors.forEach(c => c.start());
        this.jobs.forEach(j=>j.start())
        Object.keys(this.userJobs).forEach(k => this.userJobs[k].start())
    }

    stop() {
        log.info({id: this.model._id}, '停止设备');
        this.collectors.forEach(c => c.cancel());
        this.jobs.forEach(j=>j.cancel())
        Object.keys(this.userJobs).forEach(k => this.userJobs[k].cancel())
    }

    open(tunnel, model) {
        log.info({sn: tunnel.sn, id: model._id}, '打开设备');

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
                const validator = new Validator(v, this.validators)
                validator.on('alarm', alarm => {
                    Object.assign(alarm, {
                        device_id: this.model._id,
                        device_name: this.model.name || this.element.name, //添加设备名
                    });
                    log.info(alarm, '设备告警');
                    this.emit('alarm', alarm);
                });
                validator.on('error', err => {
                    this.emit('error', err);
                })
                this.validators.push(validator);
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

                const job = new Job(j, this.variables);
                job.on('error', err => this.emit('error', err));
                this.jobs.push(job);
            })

            //创建采集器
            this.adapter = new Adapter(tunnel.protocol, model.slave, element.data_points)
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
        log.info({id: this.model._id}, '关闭设备');

        //记录关闭
        createEvent({device_id: this.model._id, event: '关闭'})

        //不再监听关闭事件
        this.tunnel.off('close', this.onTunnelClose)
        this.tunnel.off('online', this.onTunnelOnline);
        this.tunnel.off('offline', this.onTunnelOffline);

        this.closed = true;
        //this.variables = {$online: false};//清空变量
        this.variables.$online = false;
        this.collectors.forEach(c => c.cancel());
        this.collectors = [];
        this.validators.forEach(v => v.cancel());
        this.validators = []; //检查器中有历史信息
        this.jobs.forEach(j => j.cancel());
        this.jobs = [];
        this.scripts = [];

        this.clearUserJob();
    }

    addCollector(c) {
        const col = new Collector(this.adapter, c);

        col.on("error", err => {
            this.emit('error', err);
        })

        col.on("data", values => {
            this.variables.$last = new Date();

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

            //合法检查（包括定时）
            this.validators.forEach(v => {
                if (!v.model.interval || !v.model.crontab)
                    v.execute()
            });

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
                        this.emit('error', err);
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
            this.emit('error', err);
        });
    }

    addUserJob(j) {
        const job = new UserJob(j, this.variables);
        if (this.userJobs[j._id])
            this.userJobs[j._id].cancel();
        this.userJobs[j._id] = job;

        job.on('error', err=>this.emit('error', err));
        job.on('execute', (command, params)=>{
            this.execute(command, params).then(() => {
                createEvent({device_id: this.model._id, event: '定时任务：' + command})
            }).catch(err => {
                createEvent({device_id: this.model._id, event: '定时任务：' + command + ' 失败：' + err.message})
            });
        })
    }

    clearUserJob() {
        Object.keys(this.userJobs).forEach(k => this.userJobs[k].cancel())
        this.userJobs = {};
    }

    removeUserJob(_id) {
        if (this.userJobs[_id]) {
            this.userJobs[_id].cancel();
            delete this.userJobs[_id];
        }
    }

    async writeValues(values) {
        for (let k in values) {
            if (!values.hasOwnProperty(k)) return;
            await this.adapter.set(k, values[k])
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
