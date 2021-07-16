const EventEmitter = require('events');
const net = require('net');
const Tunnel = require('../lib/tunnel');

module.exports = class TcpServer extends EventEmitter {
    options = {
        timeout: 30000,
        port: 0,
    }

    server;
    clients = {};

    constructor(options) {
        super()

        Object.assign(this.options, options)

        this.server = net.createServer((socket => {
            //设置超时
            socket.setTimeout(options.timeout, function () {
                console.log("连接超时了", socket.remoteAddress);
                socket.destroy()
            });

            //告诉外部，有新连接
            const tunnel = new Tunnel(socket, options.tunnel);
            this.emit('connect', tunnel)
        }));

        this.server.on("error", err => {
            //console.error(err)
            this.emit("error", err)
        })

        this.server.on("close", () => {
            this.emit("close")
        })

        this.server.listen(this.options.port);
    }

    close() {
        this.server.close();
    }
}