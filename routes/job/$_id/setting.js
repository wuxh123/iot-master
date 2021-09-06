const curd = require_plugin("mongodb/curd");
const project = require("../../../lib/project")
const device = require("../../../lib/device")

exports.post = curd.setting("job", {
    after: ctx => {
        const model = ctx.body.data;
        if (model.enable) {
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
        }
    }
});