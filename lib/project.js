const EventEmitter = require("events");
const vm = require("vm");
const _ = require("lodash");
const cron = require("./cron");
const device = require("./device");

const mongo = require_plugin("mongodb")

class Project extends EventEmitter {

    model = {};

    devices = [];

    validators = [];
    jobs = [];
    scripts = [];

    commands = {};

    //变量集，要使用defineProperty来定义变量
    variables = {};

    closed = false;


    constructor(model) {
        super();

        this.model = model;

        this.open(model);
    }

    open(model) {
        if (!this.closed) {
            this.close();
        }
        this.closed = false;

        //构建变量
        //this.variables = {};
        model.variables.forEach(v => {
            this.addVariable(v);
        })

        model.devices.forEach(d => {
            this.addDevice(d.device_id);
        })

        model.validators.forEach(v => {
            this.addValidator(v);
        })

        model.commands.forEach(c => {
            this.addCommand(c);
        })

        model.scripts.forEach(s => {
            if (s.enable) this.addScript(s);
        })

        model.jobs.forEach(j => {
            if (j.enable) this.addJob(j);
        })

        //启动用户定时任务
        mongo.db.collection("job").find({project_id: this.model._id, enable: true}).toArray().then(jobs => {
            jobs.forEach(j => this.addUserJob(j));
        }).catch(console.error);
    }

    close() {
        this.closed = true;
        this.jobs.forEach(j => j.handler.cancel());
        this.jobs = [];
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
                }
            });
            this.variables[name] = dev.variables;

            //变量
            dev.on('alarm', v => {
                //TODO 获取设备名称
                this._notice(v).then()
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
                    s.script.runInNewContext(this.variables);
                })

            })
        }
    }

    addValidator(v) {
        this.validators.push({
            script: new vm.Script(v.expression),
            model: v,
            start: 0,
            delay: v.delay * 1000 //换算成s，避免后续冗余的乘法计算
        })
    }

    addCommand(c) {
        this.commands.push({
            script: new vm.Script(c.script),
            model: c,
        })
    }

    addJob(j) {
        const job = {
            script: new vm.Script(j.script),
            model: j,
        };
        this.jobs.push(job)

        //启动定时
        job.handler = cron.schedule(j.crontab, () => {
            job.script.runInNewContext(this.variables);
        })
    }

    addUserJob(j) {
        const job = {
            model: j,
        };
        this.jobs.push(job);

        //预编译参数数组
        if (j.params)
            job.script = new vm.Script(j.params)

        //启动定时
        job.handler = cron.scheduleTime(j.time, () => {
            let params = [];
            if (j.params)
                params = job.script.runInNewContext(Object.assign({}, this.variables)); //隔离，避免操作（项目应该深度clone）
            this.execute(j.command, params)
        })
    }

    addScript(c) {
        this.scripts.push({
            script: new vm.Script(c.script),
            model: c,
        });
    }

    validate() {
        this.validators.forEach(v => {
            //const ret = vm.runInNewContext(v.expression, this.variables);
            const ret = v.script.runInNewContext(this.variables);
            if (ret) {
                v.start = 0;//去掉发生时间，重置延时
                v.reported = false;//去掉已经上报标识
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

            }).catch(console.error)
        })
    }


    execute(cmd, param) {
        const command = this.commands[cmd];
        //以下代码会污染variables，具体再定
        param.forEach((v, i) => this.variables['$' + (i + 1)] = v);
        command.script.runInNewContext(this.variables);
        //已验证 Object.assign({}, this.variables, params)作为context，set失效
        //每次都构造context，可能太低效了

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
        const pipeline = [{$match: Object.assign({}, {project_id: this.model._id}, filter)}];
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

        const subs = await mongo.db.collection("subscribe").aggregate(pipeline).toArray();
        subs.forEach(sub => {
            //TODO 语音通知需要等待结果，不成功发下一条，并做记录
            console.log('通知：', sub)
        });
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
