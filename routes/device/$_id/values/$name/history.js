const influx = require_plugin("influxdb");
const mongo = require_plugin("mongodb");

exports.post = async ctx => {
    const dvc = await mongo.db.collection("device").findOne({_id: ctx.params._id});

    const body = ctx.request.body || {};
    const res = await influx.query(dvc.element_id, {id: ctx.params._id.toString()}, ctx.params.name, body.window || '10m', body.start || '-5h', body.end || '0h')

    ctx.body = {data: res}
}