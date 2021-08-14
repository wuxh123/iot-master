const prj = require("../../../lib/project");

exports.get = async ctx=>{
    const p = prj.get(ctx.params._id);
    if (p) {
        p.close();
    }
}