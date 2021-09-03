const tunnel = require("../../../lib/tunnel");

exports.get = async ctx =>{
    const tnl = tunnel.get(ctx.params._id)
    if (!tnl) {
        ctx.body = {data: "通道未上线"}
        return
    }

    // check if the current request is websocket
    if (ctx.ws) {
        const ws = await ctx.ws() // retrieve socket
        // ws.on('message', function(message) {
        //     //console.log('received: %s', message);
        // });

        //开始透传
        tnl.transfer(ws);

        // now you have a ws instance, you can use it as you see fit
        return ws.send(JSON.stringify({type:"connected", data: ctx.params._id}))
    }

    // we're back to regular old http here
    ctx.body = '请使用WebSocket协议'
}