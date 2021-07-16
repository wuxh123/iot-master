const EventEmitter = require('events');

class Tunnel extends EventEmitter {
    conn;
    filters = [];
    adapter;

    options = {
        filters: []
    };

    constructor(conn, options) {
        super();
        this.conn = conn;
        Object.assign(this.options, options);

        this.options.filters.forEach(f=>{

            //从acceptors中找到脚本
            const filter = new require('../acceptors/filters/' + f.type)(this, f.options);
            this.filters.push(filter);
        });

        if (this.options.adapter)
            this.adapter = new require('../adapters/' + this.options.adapter.type)(this, this.options.adapter.options);

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

    write(data) {
        this.emit('write', data);
        return this.conn.write(data);
    }

    close() {
        this.conn.destroy();
    }

    _handle(data) {
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

class Acceptor extends EventEmitter {
    acceptor;

    constructor(type, options) {
        super();

        //从acceptors中找到脚本
        this.acceptor = new require('../acceptors/' + type)(options)

        this.acceptor.on('connect', conn => {
            const tunnel = new Tunnel(conn, options.tunnel);
            this.emit('connect', tunnel)
        })
    }

    close() {
        this.acceptor.close();
    }
}

