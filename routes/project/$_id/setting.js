const curd = require_plugin("curd");
const dvc = require("../../../lib/project");

exports.post = curd.setting("project", {
    after: ctx=>{
        //直接重启了。。。
        dvc.remove(ctx.params._id);
        if (ctx.body.data.enable)
            dvc.create(ctx.body.data)
    }
});