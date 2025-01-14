const device = require("../../lib/device");

const curd = require_plugin("mongodb/curd");
exports.post = curd.list("device", {
    joins: [{
        from: 'element',
        fields: ['name', 'image']
    }, {
        from: 'tunnel',
        fields: ['name', 'sn', 'last']
    }, {
        from: 'project',
        local: '_id',
        foreign: 'devices.device_id',
        noUnwind: true,
    }],
    after: async ctx => {
        ctx.body.data.forEach(d => {
            if (d.project)
                d.project = d.project.map(p => p.name);
        })
    },
});