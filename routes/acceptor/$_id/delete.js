const {removeAcceptor} = require("../../../lib/acceptor");
const curd = require_plugin("mongodb/curd");
exports.delete = exports.get = curd.delete("acceptor", {
    after: ctx=>{
        removeAcceptor(ctx.params._id)
    }
});