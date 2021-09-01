const EventEmitter = require('events');
const filter = require('./filter');
const protocol = require('./protocol');
const internal = require("./interval");
const _ = require("lodash");
const {createEvent} = require("./event");

const mongo = require_plugin("mongodb");

class Tunnel extends EventEmitter {

    /**
     * 数据模型
     * @type Object
     */
    model = {
        filters: [],
        timeout: 50,
    };

    /**
     * 连接
     * @type Socket
     */
    conn;

    /**
     * 过滤器
     * @type {*[]}
     */
    filters = [];

    /**
     * 适配器
     */
    protocol;

    /**
     * 透传通道
     * @type Socket
     */
    traversal;

    //超时检查
    checker;

    online = false;

    onSocketData = data => {
        this.last = Date.now();

        this.emit('read', data);
        this._handle(data);
    };
    onSocketError = err => {
        this.emit('error', err)
    };
    onSocketClose = () => {
        this.close();
        this.emit('offline');
        createEvent({tunnel_id: this.model._id, event: '下线'});
    };

    /**
     * 初始化
     * @param {Socket} conn 网络连接
     * @param {Object} model 参数
     */
    constructor(conn, model) {
        super();

        this.open(conn, model);


        this.on('control', val => {
            mongo.db.collection("tunnel").updateOne({_id: this.model._id}, {$set: val}).then(() => {
            }).catch(err => this.emit('error', err));
        })

        //避免无监听导致异常退出
        this.on('error', log.error);
    }

    open(conn, model) {
        log.info({sn:model.sn, id:model._id}, '打开通道');

        if (this.online) {
            log.info({sn: model.sn}, '重新上线')
            this.close()
        }

        this.conn = conn;
        _.extend(this.model, model)

        this.emit('online');
        this.online = true;
        createEvent({tunnel_id: this.model._id, event: '上线'});

        //超时
        const now = Date.now();
        const offline = now - (this.last || 0);
        this.last = now;
        if (this.model.timeout) {
            this.checker = internal.check(1000, (now) => {
                const tm = now - this.last;
                if (tm > this.model.timeout * 1000) {
                    conn.destroy()
                    //this.close();
                    this.emit('error', new Error('超时拆线' + model.sn + ' ' + tm))
                }
            });
        }

        //创建过滤器
        this.filters = [];
        if (this.model.heartbeat && this.model.heartbeat.enable)
            this.filters.push(filter.create(this, 'heartbeat', this.model.heartbeat))
        if (this.model.control && this.model.control.enable) {
            const control = filter.create(this, this.model.control.type, this.model.control.options);
            this.filters.push(control);

            let variables = ['rssi'];
            this.model.control.variables && this.model.control.variables.forEach(v => {
                //如果没有，就采集一遍
                if (!this.model.hasOwnProperty(v))
                    variables.push(v);
            })

            //这段代码好搓~~~
            variables.forEach((v, i) => setTimeout(() => control.query(v), 10000 + i * 5000));
        }

        //创建适配器
        if (this.model.protocol && this.model.protocol.type) {
            if (this.protocol) {
                this.protocol.open(this);
            } else {
                this.protocol = protocol.create(this, this.model.protocol);
            }
        }

        //监听事件
        conn.on('data', this.onSocketData);
        conn.on('error', this.onSocketError);
        conn.on('close', this.onSocketClose);
    }

    /**
     * 关闭通道
     */
    close() {
        log.info({sn:this.model.sn, id:this.model._id}, '关闭通道');

        this.conn.destroy();
        this.conn.off('data', this.onSocketData);
        this.conn.off('error', this.onSocketError);
        this.conn.off('close', this.onSocketClose);
        this.conn = undefined;
        this.online = false;
        if (this.checker)
            this.checker.cancel();
        this.checker = undefined;
    }

    /**
     * 开启透传
     * socket关闭，即自动关闭
     * 不支持断线重连
     * @param {Socket} socket
     */
    traverse(socket) {
        this.traversal = socket;
        socket.on('data', data => this.conn.write(data));
        socket.on('error', log.error);
        socket.on('close', () => this.traversal = undefined);
    }


    /**
     * 发送数据
     * @param {Buffer|string} data
     * @returns boolean
     */
    write(data) {
        if (!this.conn)
            throw new Error("通道离线");

        log.trace({sn: this.model.sn, data}, '通道写入');

        //透传过程中，不允许发数据
        if (this.traversal) {
            return false;
        }

        this.emit('write', data);
        return this.conn.write(data);
    }


    _handle(data) {
        log.trace({sn: this.model.sn, data}, '通道接收');

        let that = this;
        let index = 0;

        function next(data) {
            if (index === that.filters.length) {

                //在这里透传，可以滤掉注册包，心跳包等 非设备相关的内容
                if (that.traversal) {
                    that.traversal.write(data);
                    return;
                }

                //发送协议解析
                // if (that.adapter) {
                //     that.adapter.handle(data);
                //     return
                // }
                that.emit('data', data);
            }
            const filter = that.filters[index++];
            filter && filter.handle(data, next);
        }

        //递归执行过滤器
        next(data);
    }
}

const tunnels = {};


/**
 * 打开通知（或恢复）
 * @param {Socket} socket
 * @param {Object} model
 * @returns {Tunnel}
 */
exports.open = function (socket, model) {
    //离线恢复 逻辑
    let tunnel = tunnels[model._id]
    if (tunnel) {
        tunnel.open(socket, model);
    } else {
        tunnel = new Tunnel(socket, model);
        //建立索引
        tunnels[model._id] = tunnel;
    }

    return tunnel;
}

/**
 * 获取通知
 * @param {string} id
 * @returns {Tunnel}
 */
exports.get = function (id) {
    return tunnels[id];
}

/**
 * 删除通道
 * @param id
 */
exports.remove = function (id) {
    const tunnel = tunnels[id];
    if (tunnel) {
        tunnel.close();
        delete tunnels[id];
    }
}


