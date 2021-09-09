const device = require("../../../../lib/device");

exports.get = async ctx => {
    const dvc = device.get(ctx.params._id);
    if (!dvc) throw new Error("未上线");
    ctx.body = {data: dvc.context.values()};
}