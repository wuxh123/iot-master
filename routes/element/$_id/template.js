const curd = require_plugin("mongodb/curd");
exports.post = curd.list("template", {
    before: ctx => {
        const body = ctx.request.body;
        body.filter.devices = {$elemMatch: {element_id: ctx.params._id}};
    },
});