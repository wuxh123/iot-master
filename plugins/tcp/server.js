const EventEmitter = require('events');
const net = require('net');

module.exports = class Server extends EventEmitter {
    model = {
        timeout: 30000,
        port: 0,
    }

    register = { regex: /^\w+$/ };

    server;
    clients = {};
    closed = true;
    error = '';

    constructor(model) {
        super()
        this.open(model)
    }

    open(model) {
        if (model)
            this.model = model;
        if (!this.closed)
            this.close();
        //this.closed = false;

        if (this.model.register.regex)
            this.register.regex = new RegExp(this.model.register.regex);

        this.server = net.createServer((socket => {
            //设置超时
            if (this.model.timeout)
                socket.setTimeout(this.model.timeout * 1000, function () {
                    console.log("连接超时了", socket.remoteAddress);
                    socket.destroy()
                });

            //接收注册包
            socket.once('data', data => {
                const sn = data.toString();
                if (this.register.regex) {
                    if (!this.register.regex.test(sn)) {
                        socket.end("invalid sn")
                        return;
                    }
                }

                //告诉外部，有新连接
                //const tunnel = new Tunnel(socket, this.model);
                this.emit('connect', sn, socket);
            })

        }));

        this.server.on("error", err => {
            //console.error('server', err)
            this.error = err.message;
            this.emit("error", err);
        })

        this.server.on("close", () => {
            this.emit("close")
        })

        this.server.listen(this.model.port, () => {
            this.closed = false;
        });
    }

    close() {
        if (!this.closed) {
            if (this.server) {
                this.server.close();
                this.server = null;
            }
        }
        this.closed = true;
    }
}
