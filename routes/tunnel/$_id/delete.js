const {removeTunnel} = require("../../../lib/acceptor");
const curd = require_plugin("curd");
exports.delete = exports.get = curd.delete("tunnel", {
    after: ctx=>{
        removeTunnel(ctx.params._id)
    }
});