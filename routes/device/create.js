const curd = require_plugin("mongodb/curd");
const mongo = require_plugin("mongodb");
const dvc = require("../../lib/device")
const tnl = require("../../lib/tunnel")

exports.post = curd.create("device", {
    after: ctx=>{
        mongo.db.collection("device").findOne({_id: ctx.body.data}).then(model=>{
            const tunnel = tnl.get(model.tunnel_id);
            if (tunnel)
                dvc.open(tunnel, model)
        }).catch()
    }
});