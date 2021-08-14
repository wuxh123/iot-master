const prj = require("../../../lib/project");
const mongo = require_plugin("mongodb");

exports.get = async ctx=>{
    const p = await mongo.db.collection("project").findOne({_id: ctx.params._id});
    if (p)
        prj.create(p);
}