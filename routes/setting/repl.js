const repl = require('repl');
const websocket = require('ws');

exports.get = async ctx => {
    if (!ctx.ws) {
        ctx.body = '请使用WebSocket协议'
        return
    }

    const ws = await ctx.ws()
    const stream = websocket.createWebSocketStream(ws, {encoding: 'utf8'});

    repl.start({
        //prompt: '> ',
        input: stream,
        output: stream
    }).on('exit', () => {
        ws.close();
    });

    ws.on('error', err=>{
        console.error('ws err', err)
    })

}