const mongo = require_plugin("mongodb");

exports.post = async ctx => {
    const body = ctx.request.body;

    const stages = mongo.bodyToStages(ctx.request.body);

    //ctx.request.body
    const res = await mongo.collection("acceptor").aggregate(stages).toArray();
    ctx.body = {data: res}
}