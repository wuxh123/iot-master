const device = require("../../../lib/device");
const _ = require("lodash");
const curd = require_plugin("mongodb/curd");
exports.post = curd.list("tunnel", {
    before: ctx => {
        ctx.state.stages = [
            {$match: {company_id: ctx.params._id}},
            {
                $lookup: {
                    from: 'device',
                    localField: '_id',
                    foreignField: 'tunnel_id',
                    as: 'device'
                }
            },
            {$unwind: {path: '$device'}},
            {$replaceRoot: {newRoot: '$device'}},
        ];
    },
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
            d.project = d.project.map(p=>p.name);
            const dvc = device.get(d._id);
            if (dvc) {
                d.online = true;
                d.closed = dvc.closed;
                d.error = dvc.error;
            }
        })
    },
});