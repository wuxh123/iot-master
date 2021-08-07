const EventEmitter = require("events");
const collector = require("./collector");
const vm = require("vm");
const _ = require("lodash");
const cron = require("./cron");

const influx = require_plugin("influxdb")
const mongo = require_plugin("mongodb")

class Device extends EventEmitter {
    /**
     * @type Adapter
     */
    adapter;

    model = {};

    collectors = [];
    validators = [];
    jobs = [];
    scripts = [];

    commands = {};

    //变量集，要使用defineProperty来定义变量
    variables = {};

    addVariable(v) {
        let val = eval(v.default); //计算默认值
        Object.defineProperty(this.variables, v.name, {
            get() {
                return val;
            },
            set(value) {
                val = value;
                //需要类型，转换
                const data = this.adapter.buildData(v.type, value);
                this.adapter.write(this.model.slave, v.code, v.address, data)
                    .then(console.log).catch(console.error);
                //TODO 报出设备错误，写日志
            }
        })
    }

    addCollector(c) {
        const col = collector.create(this.adapter, this.model.slave, c);
        col.on("data", data => {

            //解析数据
            const values = this.adapter.parseData(this.model.variables, data);

            //找出要保存的
            const stores = this.model.variables.filter(v => (v.name in values) && v.store).map(v => {
                return {
                    name: v.name,
                    type: v.type,
                    value: values[v.name]
                }
            });

            //找出变化的 key
            const modify = Object.keys(values).filter(k => this.values[k] !== values[k]);

            //统一赋值
            Object.assign(this.variables, values);

            //存入数据库
            if (stores.length) {
                //TODO 改用适配器，以兼容其他时序数据库
                const table = this.model.table || this.model.type || this.model.name;
                influx.write(table, [{name: 'id', value: this.model.value}], values)
            }

            //合法检查
            this.validate();

            //向project汇报
            this.emit("values", values);

            //监听变化
            this.model.scripts.forEach(s => {
                if (s.watches && s.watches.length && !_.intersection(s.watches, modify).length) return;
                vm.runInNewContext(s.script, this.variables);
                //TODO new vm.Script(s.script) 可以先编译代码提速
            })

        });
        this.collectors.push(col);
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

            this.emit('alarm', v.content);
            //TODO content是不是也可以使用表达式（脚本），但是变量要先拷贝，避免被覆盖
        })
    }

    constructor(adapter, model) {
        super();
        this.adapter = adapter;
        Object.assign(this.model, model);

        //构建变量
        this.model.variables.forEach(v => {
            this.addVariable(v);
        })

        this.model.validators.forEach(v => {
            this.addValidator(v)
        })

        this.model.collectors.forEach(c => {
            if (c.enable) {
                const col = collector.create(adapter, model.slave, c);
                this.collectors.push(col);
            }
        })

        this.model.commands.forEach(c => {
            this.addCommand(c);
        })

        this.model.scripts.forEach(s => {
            if (s.enable)
                this.addScript(s)
        })

        this.model.jobs.forEach(s => {
            if (s.enable)
                this.addScript(s)
        })

        //启动用户定时任务
        mongo.db.collection("job").find({device_id: this.model._id, enable: true}).toArray().then(jobs => {
            jobs.forEach(j => this.addUserJob(j));
        }).catch(console.error)

    }

    close() {
        this.collectors.forEach(c => c.cancel());
        this.jobs.forEach(j => j.handler.cancel());
    }

    execute(cmd, param) {
        const command = this.commands[cmd];
        //以下代码会污染variables，具体再定
        param.forEach((v, i) => this.variables['$' + (i + 1)] = v);
        command.script.runInNewContext(this.variables);
        //已验证 Object.assign({}, this.variables, params)作为context，set失效
        //每次都构造context，可能太低效了
    }

}

exports.create = function (adapter, model) {

}