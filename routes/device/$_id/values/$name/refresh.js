const device = require("../../../../../lib/device");

exports.get = async ctx => {
    const d = device.get(ctx.params._id);
    if (!d) throw new Error("设备未上线");
    const val = await d.adapter.get(ctx.params.name);
    ctx.body = {data: val};
}