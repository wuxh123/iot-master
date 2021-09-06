const curd = require_plugin("mongodb/curd");
const dvc = require("../../../lib/device")

exports.delete = exports.get = curd.delete("device", {
    after: ctx=>{
        dvc.remove(ctx.params._id)
    }
});