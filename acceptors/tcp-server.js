const EventEmitter = require('events');
const net = require('net');

class TcpServer extends EventEmitter {
    model = {
        timeout: 30000,
        port: 0,
    }

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
                if (this.model.register.regex) {
                    if (!this.model.register.regex.test(sn)) {
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

module.exports = function (model) {
    return new TcpServer(model);
}