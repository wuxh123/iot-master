const curd = require_plugin("mongodb/curd");
exports.post = curd.list("project", {
    before: ctx => {
        const body = ctx.request.body;
        body.filter.devices = {$elemMatch: {device_id: ctx.params._id}};
    },
    joins: [{
        from: 'template',
        fields: ['name']
    }]
});