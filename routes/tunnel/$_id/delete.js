const tunnel = require("../../../lib/tunnel");
const curd = require_plugin("mongodb/curd");
exports.delete = exports.get = curd.delete("tunnel", {
    after: ctx=>{
        tunnel.remove(ctx.params._id)
    }
});