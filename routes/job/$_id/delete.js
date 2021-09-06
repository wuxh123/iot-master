const curd = require_plugin("mongodb/curd");
const mongo = require_plugin("mongodb");
const project = require("../../../lib/project")
const device = require("../../../lib/device")

exports.delete = exports.get = curd.delete("job", {
    after: ctx => {
        mongo.db.collection("job_deleted").findOne({job_id: ctx.params._id}).then(model => {
            if (model.project_id) {
                const prj = project.get(model.project_id);
                if (prj) {
                    prj.removeUserJob(ctx.params._id)
                }
            } else if (model.device_id) {
                const dvc = device.get(model.device_id);
                if (dvc) {
                    dvc.removeUserJob(ctx.params._id)
                }
            }
        }).catch()
    }
});