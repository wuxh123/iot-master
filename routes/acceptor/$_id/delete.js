const {removeAcceptor} = require("../../../lib/acceptor");
const curd = require_plugin("curd");
exports.delete = exports.get = curd.delete("acceptor", {
    after: ctx=>{
        removeAcceptor(ctx.params._id)
    }
});