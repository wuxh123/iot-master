const EventEmitter = require('events');
const filter = require('./filter');
const adapter = require('./adapter');
const internal = require("./interval");
const _ = require("lodash");

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
    adapter;

    /**
     * 透传通道
     * @type Socket
     */
    traversal;

    //超时检查
    checker;


    onSocketData = data=>{
        this.last = Date.now();

        this.emit('read', data);
        this._handle(data);
    };
    onSocketError = err=>{
        this.emit('error', err)
    };
    onSocketClose = () => {
        this.conn.off('close', this.onSocketClose);
        this.checker.cancel();
        this.emit('offline');

    };

    /**
     * 初始化
     * @param {Socket} conn 网络连接
     * @param {Object} model 参数
     */
    constructor(conn, model) {
        super();

        this.open(conn, model)
    }

    open(conn, model) {
        this.conn = conn;
        _.extend(this.model, model)

        //超时
        this.last = Date.now();
        this.checker = internal.check('2000', (now)=>{
            if (this.last + this.model.timeout * 1000 < now) {
                conn.destroy()
            }
        });

        //创建过滤器
        this.filters = [];
        if (this.model.heartbeat.enable)
            this.filters.push(filter.create(this, 'heartbeat', this.model.heartbeat))
        if (this.model.control.enable)
            this.filters.push(filter.create(this, 'control', this.model.control))

        //创建适配器
        if (this.model.adapter && this.model.adapter.type) {
            if (this.adapter) {
                this.adapter.attach(this, this.model.adapter);
            } else {
                this.adapter = adapter.create(this, this.model.adapter);
            }
        }

        //监听事件
        this.conn.on('data', this.onSocketData);
        this.conn.on('error', this.onSocketError);
        this.conn.on('close', this.onSocketClose);
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
        socket.on('error', console.error);
        socket.on('close', () => this.traversal = undefined);
    }


    /**
     * 发送数据
     * @param {Buffer|string} data
     * @returns boolean
     */
    write(data) {
        //透传过程中，不允许发数据
        if (this.traversal) {
            return false;
        }

        this.emit('write', data);
        return this.conn.write(data);
    }

    /**
     * 关闭通道
     */
    close() {
        this.conn.destroy();
    }

    _handle(data) {

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


exports.open = function (socket, model) {
    //离线恢复 逻辑
    let tunnel = tunnels[model._id]
    if (tunnel) {
        tunnel.open(socket, model);
    } else {
        tunnel = new Tunnel(tunnel, model);
        //建立索引
        tunnels[model._id] = tunnel;
    }


}

exports.get = function (id) {
    return tunnels[id];
}

exports.remove = function (id) {
    const tunnel = tunnels[id];
    if (tunnel) {
        tunnel.close();
        delete tunnels[id];
    }
}


