const curd = require_plugin("curd");
exports.post = curd.list("member", {
    before: ctx=>{
        const body = ctx.request.body;
        body.filter.user_id = ctx.params._id;
    },
    join:{
        from: 'group',
        replace: true
    }
});