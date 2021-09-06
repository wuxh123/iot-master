const acc = require("../../../lib/acceptor");
const curd = require_plugin("mongodb/curd");
exports.post = curd.setting("acceptor", {
    after: ctx=>{
        //直接重启了。。。
        acc.removeAcceptor(ctx.params._id);
        if (ctx.body.data.enable)
            acc.create(ctx.body.data)
    }
});