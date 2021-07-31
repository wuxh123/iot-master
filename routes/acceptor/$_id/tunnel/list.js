const curd = require_plugin("curd");
exports.post = curd.list("tunnel", {
    name: 'string',
    type: 'string',
    remote: 'string',
}, ctx => {
    //强制修改为
    const filter = {key: "acceptor_id", value: ctx.params._id};
    ctx.body.filter = ctx.body.filter || [];
    ctx.body.filter.push(filter);
});