const os = require('os');
const { spawn } = require('child_process');
const websocket = require('ws');

exports.get = async ctx => {
    if (!ctx.ws) {
        ctx.body = '请使用WebSocket协议'
        return
    }

    const ws = await ctx.ws();
    const stream = websocket.createWebSocketStream(ws, {encoding: 'utf8'});

    const sh = spawn(os.platform() === 'win32' ? 'cmd' : 'bash');

    sh.stdout.pipe(stream);
    sh.stderr.pipe(stream);
    stream.pipe(sh.stdin);
}