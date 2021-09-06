const curd = require_plugin("mongodb/curd");
exports.post = curd.list("member", {
    before: ctx=>{
        const body = ctx.request.body;
        body.filter = body.filter || {};
        body.filter.company_id = ctx.params._id;
    },
    join:{
        from: 'user'
    }
});