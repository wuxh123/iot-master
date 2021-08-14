const curd = require_plugin("curd");
const dvc = require("../../../lib/device");
const acc = require("../../../lib/acceptor");

exports.post = curd.setting("device", {
    after: ctx=>{
        //直接重启了。。。
        dvc.remove(ctx.params._id);
        if (ctx.body.data.enable) {
            const tunnel = acc.getTunnel(ctx.body.data.tunnel_id);
            if (tunnel)
                dvc.create(tunnel, ctx.body.data)
        }
    }
});