const curd = require_plugin("curd");
exports.post = curd.list("tunnel", {
    name: 'string',
    type: 'string',
    remote: 'string',
}, ctx => {
    //强制修改为
    const body = ctx.request.body;
    const filter = {key: "acceptor_id", value: ctx.params._id};
    body.filter = body.filter || [];
    body.filter.push(filter);
});