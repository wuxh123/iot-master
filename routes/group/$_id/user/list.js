const curd = require_plugin("curd");
exports.post = curd.list("user", {
    name: 'string',
}, ctx => {
    //强制修改为
    const filter = {key: "group_id", value: ctx.params._id};
    ctx.body.filter = ctx.body.filter || [];
    ctx.body.filter.push(filter);
});