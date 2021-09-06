const curd = require_plugin("mongodb/curd");
const mongo = require_plugin("mongodb");
const prj = require("../../lib/project")

exports.post = curd.create("project", {
    after: ctx=>{
        mongo.db.collection("project").findOne({_id: ctx.body.data}).then(model=>{
            prj.create(model)
        }).catch()
    }
});