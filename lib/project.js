const EventEmitter = require("events");
const _ = require("lodash");
const cron = require("./cron");
const device = require("./device");
const script = require("./script");
const {notice} = require("./notice");
const {createEvent} = require("./event");

const mongo = require_plugin("mongodb")

class Selector {
    devices = [];

    constructor(dvcs) {
        this.devices = dvcs;
    }

    exec(cmd, params) {
        this.devices.forEach(d => {
            createEvent({device_id: d.model._id, event: '批量执行：' + cmd})
            d.execute(cmd, params).then(()=>{}).catch(log.error);
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
        return this.devices[0].variables;
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
    variables = {};

    closed = true;
    error = '';


    constructor(model) {
        super();

        const that = this;
        this.variables.$ = function () {
            //console.log(model.devices)
            const tags = arguments;
            const dvcs = [];
            for (let k in that.devices) {
                if (!that.devices.hasOwnProperty(k)) return;
                const d = that.devices[k].instance;
                if (_.intersection(tags, d.element.tags).length)
                    dvcs.push(d)
            }
            return new Selector(dvcs);
        }

        this.open(model);
    }

    init(model) {
        //构建变量
        //this.variables = {};
        model.variables && model.variables.forEach(v => {
            this.addVariable(v);
        })

        model.validators && model.validators.forEach(v => {
            this.addValidator(v);
        })

        model.commands && model.commands.forEach(c => {
            this.addCommand(c);
        })

        model.scripts && model.scripts.forEach(s => {
            if (s.enable) this.addScript(s);
        })

        model.jobs && model.jobs.forEach(j => {
            if (j.enable) this.addJob(j);
        })
    }

    open(model) {
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
        this.closed = true;
        this.jobs.forEach(j => j.handler.cancel());
        this.jobs = [];
        for (let id in this.userJobs) {
            if (this.userJobs.hasOwnProperty(id))
                this.userJobs[id].handler.cancel()
        }
        this.userJobs = {};
        //删除设备，并取消消息监听
        for (let id in this.devices) {
            if (this.devices.hasOwnProperty(id))
                this.removeDevice(id)
        }

        this.clearUserJob();
    }

    addVariable(v) {
        this.variables[v.name] = eval(v.default); //计算默认值

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
                this.variables[name] = dev.variables;
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
                        s.script.runInNewContext(this.variables);
                    } catch (err) {
                        this.error = err.message;
                        log.error(err.message)
                    }
                })
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
                log.error(err.message)
            }
        })
    }

    initUserJob() {
        //启动用户定时任务
        mongo.db.collection("job").find({project_id: this.model._id, enable: true}).toArray().then(jobs => {
            jobs.forEach(j => this.addUserJob(j));
        }).catch(err => {
            this.error = err.message;
            log.error(err.message)
        });
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

                //记录定时任务
                createEvent({project_id: this.model._id, event: '定时任务：' + j.command})

                this.execute(j.command, params)
            } catch (err) {
                this.error = err.message;
                log.error(err.message)
            }
        })
    }

    clearUserJob() {
        Object.keys(this.userJobs).forEach(k => this.userJobs[k].handler.cancel())
        this.userJobs = {};
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
                log.error(err.message)
                return;
            }
            //以下是不合法处理逻辑

            //延时处理，发生时间，当前时间
            const now = Date.now() * 0.001;
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
            v.reported = true;

            //处理通知
            notice({
                name: v.name,
                content: v.content,
                level: v.level,
                project_name: this.model.name, //添加项目名
                project_id: this.model._id,
                group_id: this.model.group_id,
                company_id: this.model.company_id,
            }).then().catch(err => {
                this.error = err.message;
                log.error(err.message)
            })
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
            log.error(err.message)
        }

        //设备离线，可以缓存一段时间？
    }
}

const projects = {};

exports.create = function (model) {
    //TODO 离线恢复 逻辑

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
