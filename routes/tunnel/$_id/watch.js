const acceptor = require("../../../lib/acceptor");


exports.get = async ctx =>{
    const tunnel = acceptor.getTunnel(ctx.params._id)
    if (!tunnel) {
        ctx.body = {data: "通道未上线"}
        return
    }

    // check if the current request is websocket
    if (ctx.ws) {
        const ws = await ctx.ws() // retrieve socket
        function onRead(data){
            ws.send(JSON.stringify({
                type: 'read',
                size: data.length,
                data: data.toString(),
                hex: data.toString('hex'),
            }))
        }

        function onWrite(data){
            ws.send(JSON.stringify({
                type: 'write',
                size: data.length,
                data: data.toString(),
                hex: data.toString('hex'),
            }))
        }

        function onError(err){
            ws.send(JSON.stringify({
                type: 'error',
                data: err.message,
            }))
        }

        tunnel.on('read', onRead)

        tunnel.on('write', onWrite)

        tunnel.on('error', onError)

        ws.on('close', ()=>{
            tunnel.off('read', onRead);
            tunnel.off('write', onWrite);
            tunnel.off('error', onError);
        })

        ws.on('message', function(message) {
            //console.log('received: %s', message);
            const obj = JSON.parse(message)
            switch (obj.type) {
                case 'write':
                    const d = Buffer.from(obj.data, obj.isHex ? 'hex' : 'utf8');
                    tunnel.write(d)
                    break;
            }
            //tunnel.write()
        });

        // now you have a ws instance, you can use it as you see fit
        return ws.send('hello')
    }

    // we're back to regular old http here
    ctx.body = '请使用新浏览器'
}