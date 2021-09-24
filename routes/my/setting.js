const mongo = require_plugin("mongodb");
exports.post = (async ctx=>{
    const body = ctx.request.body;
    const ret = await mongo.db.collection("user").updateOne({_id: ctx.state.user._id}, {$set: body})
    ctx.body = {data: ret}
});
