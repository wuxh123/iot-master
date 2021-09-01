const EventEmitter = require('events');
const net = require('net');
const Tunnel = require('../lib/tunnel');

let server;

exports.start = function (opts) {
    let options = {
        timeout: 30000,
        port: 1843,
    }

    Object.assign(options, opts)

    server = net.createServer((socket => {
        //设置超时
        socket.setTimeout(options.timeout, function () {
            //console.log("连接超时了", socket.remoteAddress);
            socket.destroy()
        });

        socket.once("data", data => {
            //TODO 解析注册码

            //TODO 找到通道，进行透传


        });
    }));

    server.on("error", err => {
        //log.error(err)
    })

    server.on("close", () => {

    })

    server.listen(options.port);
}


exports.stop = function () {
    server.close();
    server = undefined;
}