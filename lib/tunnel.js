const EventEmitter = require('events');
const filter = require('./filter');
const adapter = require('./adapter');

module.exports = class Tunnel extends EventEmitter {
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
     * @type net.Socket
     */
    traversal;

    options = {
        filters: []
    };

    /**
     * 初始化
     * @param {Socket} conn 网络连接
     * @param {Object} options 参数
     */
    constructor(conn, options) {
        super();
        this.conn = conn;
        Object.assign(this.options, options);

        //创建过滤器
        if (this.options.register)
            this.filters.push(filter.create(this, this.options.register))
        if (this.options.heartbeat)
            this.filters.push(filter.create(this, this.options.heartbeat))
        if (this.options.control)
            this.filters.push(filter.create(this, this.options.control))

        //创建适配器
        if (this.options.adapter)
            this.adapter = adapter.create(this, this.options.adapter);


        this.conn.on('data', data => {
            this.emit('read', data);
            this._handle(data);
        });

        this.conn.on('error', error => {
            this.emit('error', error);
        });

        this.conn.on('close', () => {
            this.emit('close');
        });
    }

    /**
     * 开启透传
     * socket关闭，即自动关闭
     * 不支持断线重连
     * @param {net.Socket} socket
     */
    traverse(socket) {
        this.traversal = socket;
        socket.on('data', data => this.conn.write(data));
        socket.on('error', error => socket.destroy());
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
                that.adapter.handler(data);
                return
            }
            const filter = that.filters[index++];
            filter && filter.handle(data, next);
        }

        //递归执行过滤器
        next(data);
    }
}
