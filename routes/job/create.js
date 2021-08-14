const curd = require_plugin("curd");
const mongo = require_plugin("mongodb");
const project = require("../../lib/project")
const device = require("../../lib/device")

exports.post = curd.create("job", {
    after: ctx => {
        mongo.db.collection("job").findOne({_id: ctx.body.data}).then(model => {
            if (!model.enable) return;
            if (model.project_id) {
                const prj = project.get(model.project_id);
                if (prj) {
                    prj.addUserJob(model)
                }
            } else if (model.device_id) {
                const dvc = device.get(model.device_id);
                if (dvc) {
                    dvc.addUserJob(model)
                }
            }
        }).catch()
    }
});