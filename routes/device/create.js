const curd = require_plugin("curd");
const mongo = require_plugin("mongodb");
const dvc = require("../../lib/device")
const acc = require("../../lib/acceptor")

exports.post = curd.create("device", {
    after: ctx=>{
        mongo.db.collection("device").findOne({_id: ctx.body.data}).then(model=>{
            const tunnel = acc.getTunnel(model.tunnel_id);
            if (tunnel)
                dvc.create(tunnel, model)
        }).catch()
    }
});