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
        this.options.filters.forEach(f => {
            const fi = filter.create(this, f);
            this.filters.push(fi);
        });

        //创建适配器
        if (this.options.adapter) {
            this.adapter = adapter.create(this, this.options.adapter);
        }

        //TODO 创建设备 重连

        //TODO 创建项目（模板）

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
        socket.on('data', data => this.write(data));
        socket.on('error', error => socket.destroy());
        socket.on('close', () => this.traversal = undefined);
    }

    /**
     * 发送数据
     * @param {Buffer|string} data
     * @returns boolean
     */
    write(data) {
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
        //透传
        if (this.traversal) {
            this.traversal.write(data);
            return;
        }

        let that = this;
        let index = 0;

        function next(data) {
            if (index === that.filters.length) {
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
