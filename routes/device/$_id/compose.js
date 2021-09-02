const device = require("../../../lib/device");
const _ = require("lodash");

const curd = require_plugin("curd");
exports.get = curd.compose("device", {
    joins: [{
        from: 'element',
        //fields: ['name']
    }, {
        from: 'tunnel',
        //fields: ['name', 'sn', 'last']
    }, {
        from: 'project',
        local: '_id',
        foreign: 'devices.device_id', //$elemMatch没办法使用，也就没办法加fields
        noUnwind: true,
    }],
    after: async ctx => {
        const d = ctx.body.data;
        d.project = d.project.map(p=>p.name);

        const dvc = device.get(d._id);
        if (dvc) {
            d.online = dvc.context.$online;
            d.closed = dvc.closed;
            d.error = dvc.error;
            d.values = dvc.context;
        }
    },
});