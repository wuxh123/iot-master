const mongo = require_plugin("mongodb");

const crypto = require("crypto");

function md5(str) {
    return crypto.createHash("md5").update(str).digest('hex');
}

exports.post = (async ctx => {
    const body = ctx.request.body;

    const user = await mongo.db.collection("user").findOne({_id: ctx.state.user._id});
    if (user.password === md5(body.old))
        throw new Error("原密码错误");

    const ret = await mongo.db.collection("user").updateOne({_id: ctx.state.user._id}, {
        $set: {
            password: md5(body.new)
        }
    });
    ctx.body = {data: ret}
});
