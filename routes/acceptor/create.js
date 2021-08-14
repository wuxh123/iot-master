const acc = require("../../lib/acceptor")
const mongo = require_plugin("mongodb");

const curd = require_plugin("curd");
exports.post = curd.create("acceptor", {
    after: ctx=>{
        mongo.db.collection("acceptor").findOne({_id: ctx.body.data}).then(model=>{
            if (model.enable)
                acc.create(model)
        }).catch()
    }
});