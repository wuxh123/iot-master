const curd = require_plugin("curd");
const dvc = require("../../../lib/device");
const tnl = require("../../../lib/tunnel");

exports.post = curd.setting("device", {
    after: ctx=>{
        //直接重启了。。。
        dvc.remove(ctx.params._id);
        if (ctx.body.data.enable) {
            const tunnel = tnl.get(ctx.body.data.tunnel_id);
            if (tunnel)
                dvc.open(tunnel, ctx.body.data)
        }
    }
});