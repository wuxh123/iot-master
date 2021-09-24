const EventEmitter = require("events");
const _ = require("lodash");
const cron = require("./cron");
const device = require("./device");
const script = require("./script");
const {notice} = require("./notice");
const {createEvent} = require("./event");
const Validator = require("./validator");
const Job = require("./job");
const UserJob = require("./user_job");
const Context = require("./context");

const mongo = require_plugin("mongodb")

class Selector {
    devices = [];

    constructor(dvcs) {
        this.devices = dvcs;
    }

    exec(cmd, params) {
        this.devices.forEach(d => {
            createEvent({device_id: d.model._id, event: '批量执行：' + cmd})
            d.execute(cmd, params).then(() => {
            }).catch(log.error);
        })
        return this;
    }

    each(fn) {
        this.devices.forEach(d => {
            fn && fn(d)
        });
        return this;
    }

    odd() {
        const ds = this.devices.filter((v, i) => i % 2 === 0);
        return new Selector(ds);
    }

    even() {
        const ds = this.devices.filter((v, i) => i % 2 === 1);
        return new Selector(ds);
    }

    size() {
        return this.devices.length;
    }

    value() {
        if (!this.devices.length)
            return {}; //null ???
        return this.devices[0].context;
    }
}

class Project extends EventEmitter {

    model = {};

    devices = {};

    validators = [];
    jobs = [];
    scripts = [];

    commands = {};
    userJobs = {};

    //变量集，要使用defineProperty来定义变量
    context = new Context()//{};

    closed = true;
    error = '';


    constructor(model) {
        super();

        this.on('error', err => {
            this.error = err.message;
            log.error({id: this.model._id, error: err.message}, '项目错误');
        });

        this.on('alarm', alarm => {
            notice(alarm).then().catch(err => {
                this.emit('error', err);
            });
        });

        const that = this;
        this.context.$ = function () {
            //console.log(model.devices)
            const tags = [...arguments];
            const dvcs = [];
            for (let k in that.devices) {
                if (!that.devices.hasOwnProperty(k)) return;
                const d = that.devices[k].instance;
                if (_.intersection(tags, d.element.tags).length)
                    dvcs.push(d)
            }
            return new Selector(dvcs);
        }
        this.context.$type = function () {
            //console.log(model.devices)
            const types = [...arguments];
            const dvcs = [];
            for (let k in that.devices) {
                if (!that.devices.hasOwnProperty(k)) return;
                const d = that.devices[k].instance;
                if (types.indexOf(d.element.type) > -1)
                    dvcs.push(d)
            }
            return new Selector(dvcs);
        }
        this.context.$project = this;

        this.open(model);
    }

    createEvent(event) {
        createEvent({project_id: this.model._id, event: event});
    }

    update(data) {
        mongo.db.collection("project").updateOne({_id: this.model._id}, {$set: data})
            .then(res => Object.assign(this.model, data)).catch(err => this.emit('error', err));
    }

    init(model) {
        //构建变量
        //this.variables = {};
        model.variables && model.variables.forEach(v => {
            this.addVariable(v);
        })

        model.validators && model.validators.forEach(v => {
            if (!v.enable) return;
            const validator = new Validator(v, this.context)
            validator.on('alarm', alarm => {
                Object.assign(alarm, {
                    project_id: this.model._id,
                    project_name: this.model.name, //添加项目名
                    group_id: this.model.group_id,
                    company_id: this.model.company_id,
                });
                log.info(alarm, 'project alarm');
                this.emit('alarm', alarm);
            });
            validator.on('error', err => {
                this.emit('error', err);
            })
            this.validators.push(validator);
        })

        model.commands && model.commands.forEach(c => {
            this.commands[c.name] = {
                script: script.compile(c.script),
                model: c,
            };
        })

        model.scripts && model.scripts.forEach(s => {
            if (!s.enable) return;
            this.scripts.push({
                script: script.compile(s.script),
                model: s,
            });
        })

        model.jobs && model.jobs.forEach(j => {
            if (!j.enable) return;

            const job = new Job(j, this.context);
            job.on('error', err => this.emit('error', err));
            this.jobs.push(job);
        })
    }

    open(model) {
        log.info({id: model._id}, 'open project');

        //记录上线
        this.update({online: true, last: new Date()});

        if (model)
            this.model = model;

        if (!this.closed)
            this.close();
        this.closed = false;

        if (model.template_id) {
            mongo.db.collection("template").findOne({_id: model.template_id}).then(template => {
                this.init(template || model)
            }).catch(err => {
                this.error = err.message;
                log.error(err.message)
            })
        } else {
            this.init(model)
        }

        model.devices.forEach(d => {
            this.setDevice(d.device_id);
        })

        //启动用户定时任务
        this.initUserJob();

    }

    close() {
        log.info({id: this.model._id}, 'close project');

        //记录关闭
        this.createEvent("关闭")

        //记录上线
        this.update({online: false});

        this.closed = true;
        this.jobs.forEach(j => j.cancel());
        this.jobs = [];
        this.scripts = [];
        //删除设备，并取消消息监听
        for (let id in this.devices) {
            if (this.devices.hasOwnProperty(id))
                this.removeDevice(id)
        }

        this.clearUserJob();
    }

    addVariable(v) {
        this.context[v.name] = eval(v.default); //计算默认值

    }

    removeDevice(_id) {
        const d = this.devices[_id];
        if (d) {
            d.instance.off('alarm', d.onAlarm);
            d.instance.off('modify', d.onModify);
            d.instance.off('close', d.onClose);
            delete this.devices[_id];
        }
    }

    setDevice(_id, dev) {
        if (!dev) dev = device.get(_id);
        if (!dev) return;

        //如果已经添加过，就不再添加
        if (this.devices[_id]) {
            if (this.devices[_id].instance === dev)
                return;
            this.removeDevice(_id);
        }

        //先找出别名
        let name = ''
        this.model.devices.forEach(dd => {
            if (dd.device_id.equals(_id)) {
                name = dd.name;
                //dd.device = dev;
                //this.devices[_id] = dev;
                this.context[name] = dev.variables;
            }
        });

        const dd = {
            name,
            instance: dev,
            onAlarm: alarm => {
                //获取设备名称
                notice(Object.assign({}, alarm, {
                    project_name: this.model.name, //添加项目名
                    project_id: this.model._id,
                    group_id: this.model.group_id,
                    company_id: this.model.company_id,
                })).then().catch(err => {
                    this.error = err.message;
                    log.error(err.message)
                })
            },
            onModify: modify => {
                //添加前缀
                modify = modify.map(m => name + '.' + m);
                this.scripts.forEach(s => {
                    if (s.model.watches && s.model.watches.length && !_.intersection(s.model.watches, modify).length)
                        return;
                    try {
                        s.script.runInNewContext(this.context);
                    } catch (err) {
                        this.error = err.message;
                        log.error(err.message)
                    }
                })

                //合法检查（包括定时）
                this.validators.forEach(v => {
                    //if (!v.model.interval || !v.model.crontab)
                    v.execute()
                });
            },
            onClose: () => {
                //设备下线，可以做点儿什么呢？
                //当所有设备下线，项目就可以关闭了
            }
        }
        dev.on('alarm', dd.onAlarm)
        dev.on('modify', dd.onModify)
        dev.on('close', dd.onClose)

        this.devices[_id] = dd;
    }


    initUserJob() {
        //启动用户定时任务
        mongo.db.collection("job").find({project_id: this.model._id, enable: true}).toArray().then(jobs => {
            jobs.forEach(j => this.addUserJob(j));
        }).catch(err => {
            this.emit('error', err);
        });
    }

    addUserJob(j) {
        const job = new UserJob(j, this.context);
        if (this.userJobs[j._id])
            this.userJobs[j._id].cancel();
        this.userJobs[j._id] = job;

        job.on('error', err => this.emit('error', err));
        job.on('execute', (command, params) => {
            this.execute(command, params);
            this.createEvent('定时任务：' + command);
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


    execute(cmd, param) {
        if (this.closed)
            throw new Error("设备离线")

        const command = this.commands[cmd];
        if (!command)
            throw new Error("不支持的命令")

        //创建子对象，避免污染variables
        const ctx = {};
        Array.isArray(param) && param.forEach((v, i) => ctx['$' + (i + 1)] = v);
        ctx.__proto__ = this.context;
        try {
            command.script.runInNewContext(ctx);
            //已验证 Object.assign({}, this.variables, params)作为context，set失效
            //每次都构造context，可能太低效了
        } catch (err) {
            this.emit('error', err);
        }

        //设备离线，可以缓存一段时间？
    }
}

const projects = {};

exports.create = function (model) {
    let project = projects[model._id];
    if (!project) {
        project = new Project(model);
        //建立索引
        projects[model._id] = project;
    } else {
        project.open(model)
    }

    return project;
}

exports.get = function (id) {
    return projects[id];
}

exports.remove = function (id) {
    const project = projects[id];
    if (project) {
        project.close();
        delete project[id];  //projects[id]=null
    }
}
