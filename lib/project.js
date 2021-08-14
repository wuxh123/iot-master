const EventEmitter = require("events");
const _ = require("lodash");
const cron = require("./cron");
const device = require("./device");
const script = require("./script");

const mongo = require_plugin("mongodb")

class Selector {
    devices = [];

    constructor(dvcs) {
        this.devices = dvcs;
    }

    exec(cmd, params) {
        this.devices.forEach(d => {
            d.device.execute(cmd, params);
        })
        return this;
    }

    each(fn) {
        this.devices.forEach(d => {
            fn && fn(d.device)
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
}

class Project extends EventEmitter {

    model = {};

    devices = [];

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

        this.variables.$ = function select() {
            //console.log(model.devices)
            const tags = arguments;
            const dvcs = model.devices.filter(d => d.device && _.intersection(tags, d.device.element.tags).length)
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
            })
        } else {
            this.init(model)
        }

        model.devices.forEach(d => {
            this.addDevice(d.device_id);
        })


        //启动用户定时任务
        mongo.db.collection("job").find({project_id: this.model._id, enable: true}).toArray().then(jobs => {
            jobs.forEach(j => this.addUserJob(j));
        }).catch(err => {
            this.error = err.message;
        });
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
    }

    addVariable(v) {
        this.variables[v.name] = eval(v.default); //计算默认值

    }

    addDevice(_id, dev) {
        if (!dev)
            dev = device.get(_id);
        if (dev) {
            //先找出别名
            let name = ''
            this.model.devices.forEach(dd => {
                if (dd.device_id.equals(_id)) {
                    name = dd.name;
                    dd.device = dev;
                }
            });
            this.variables[name] = dev.variables;

            //变量
            dev.on('alarm', v => {
                //TODO 获取设备名称
                this._notice(v).then().catch(err => {
                    this.error = err.message;
                })
            })

            dev.on('values', v => {
                //console.log('数据更新', this.model.name, name, v, this.variables, JSON.stringify(this.variables))

            })

            dev.on('modify', modify => {
                //添加前缀
                modify = modify.map(m => name + '.' + m);
                this.scripts.forEach(s => {
                    if (s.model.watches && s.model.watches.length && !_.intersection(s.model.watches, modify).length)
                        return;
                    try {
                        s.script.runInNewContext(this.variables);
                    } catch (err) {
                        this.error = err.message;
                    }
                })

            })
        }
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
            v.reported = true;

            //处理通知
            this._notice(v.model).then(res => {

            }).catch(err => {
                this.error = err.message;
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
        }

        //设备离线，可以缓存一段时间？
    }

    async _notice(validator) {
        const date = new Date();
        const time = date.getHours() * 60 + date.getMinutes();
        const filter = {
            enable: true,
            level: {$gte: validator.level},
            start: {$gte: time},
            end: {$lte: time},
        }
        //union查询订阅，先项目，然后是分组，最后是公司
        const pipeline = []
        pipeline.push({$match: Object.assign({}, {project_id: this.model._id}, filter)});
        //查询分组订阅
        if (this.model.group_id)
            pipeline.push({
                $unionWith: {
                    coll: 'subscribe',
                    pipeline: [{$match: Object.assign({}, {group_id: this.model.group_id}, filter)}]
                }
            })
        //查询企业订阅
        if (this.model.company_id)
            pipeline.push({
                $unionWith: {
                    coll: 'subscribe',
                    pipeline: [{$match: Object.assign({}, {company_id: this.model.company_id}, filter)}]
                }
            });

        //查询用户
        pipeline.push({
            $lookup: {
                from: 'user',
                as: 'user',
                localField: 'user_id',
                foreignField: '_id'
            }
        })
        pipeline.push({$unwind: {path: '$user'}}) //, preserveNullAndEmptyArrays: true} 找不到用户 就过滤掉 或者 TODO 删除无效订阅
        pipeline.push({$replaceRoot: {newRoot: '$user'}})

        //查出所有订阅
        const subs = await mongo.db.collection("subscribe").aggregate(pipeline).toArray();

        //TODO 根据user去重

        let smsSubs = subs.filter(sub => sub.sms && sub.user.cellphone).map(sub => sub.user.cellphone);
        smsSubs = [...new Set(smsSubs)];
        console.log('短信通知', smsSubs);

        let voiceSubs = subs.filter(sub => sub.voice && sub.user.cellphone).map(sub => sub.user.cellphone);
        voiceSubs = [...new Set(voiceSubs)];
        console.log('语音通知', voiceSubs);

        let emailSubs = subs.filter(sub => sub.email && sub.user.email).map(sub => sub.user.email);
        emailSubs = [...new Set(emailSubs)];
        console.log('邮件通知', emailSubs);

        let wxSubs = subs.filter(sub => sub.weixin && sub.user.wx && sub.user.wx.official).map(sub => sub.user.wx.official.openid);
        wxSubs = [...new Set(wxSubs)];
        console.log('微信通知', wxSubs);

    }
}

const projects = {};

exports.create = function (model) {
    //TODO 离线恢复 逻辑

    const project = new Project(model);

    //建立索引
    projects[model._id] = project;

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
