const project = require("../../lib/project");

const curd = require_plugin("mongodb/curd");
exports.post = curd.list("project", {
    after: async ctx => {
        ctx.body.data.forEach(p => {
            const prj = project.get(p._id);
            if (prj) {
                p.online = true;
                p.closed = prj.closed;
                p.error = prj.error;
            }
        })
    },
    joins: [{
        from: 'template',
        fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }, {
        from: 'group',
        fields: ['name']
    }, {
        from: 'user',
        fields: ['name']
    }]
});