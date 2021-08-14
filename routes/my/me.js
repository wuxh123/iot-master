const mongo = require_plugin("mongodb");
exports.get = (async ctx => {
    const ret = await mongo.db.collection("user").findOne({_id: ctx.state.user._id});
    ctx.body = {data: ret};
});