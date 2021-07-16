const EventEmitter = require('events');
const net = require('net');

module.exports = class TcpClient extends EventEmitter {
    options = {
        timeout: 30000,
        host: '',
        port: 0,
    }

    socket;
    clients = {};

    constructor(options) {
        super()

        Object.assign(this.options, options)

        this.socket = new net.Socket();

        this.socket.on("error", err => {
            //console.error(err)
            this.emit("error", err);

            //TODO 重连
            setTimeout(()=>this.connect(), 10000);
        })

        this.socket.on("close", () => {
            this.emit("close");

        })
    }

    connect(){
        this.socket.connect(this.options.port, this.options.host, () => {
            //告诉外部，有新连接
            this.emit("connect", this.socket);
        })
    }

    close() {
        this.server.close();
    }
}