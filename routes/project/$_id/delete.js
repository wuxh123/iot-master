const curd = require_plugin("mongodb/curd");
const dvc = require("../../../lib/project")

exports.delete = exports.get = curd.delete("project", {
    after: ctx=>{
        dvc.remove(ctx.params._id)
    }
});