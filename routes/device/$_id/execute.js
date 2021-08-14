const device = require("../../../lib/device");

const mongo = require_plugin("mongodb");
exports.post = async ctx => {
    const body = ctx.request.body;
    const d = device.get(ctx.params._id);
    if (!d) throw new Error("设备未上线")
    d.execute(body.command, body.parameters)
    //记录日志
    await mongo.db.collection("event").insertOne({
        device_id: ctx.params._id,
        event: '执行：' + body.command,
        user_id: ctx.state.user._id
    })
    ctx.body = {data: '执行成功'}
}