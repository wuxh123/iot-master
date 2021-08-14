const project = require("../../../lib/project");
const device = require("../../../lib/device");
const _ = require("lodash");

const mongo = require_plugin("mongodb");
const curd = require_plugin("curd");
exports.get = curd.compose("project", {
    joins: [{
        from: 'template',
        //fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }, {
        from: 'group',
        fields: ['name']
    }, {
        from: 'user',
        fields: ['name']
    }, {
        from: 'device',
        local: 'devices.device_id',
        noUnwind: true,
    }, {
        from: 'element',
        local: 'device.element_id',
        noUnwind: true,
    }],
    after: async ctx => {
        const p = ctx.body.data;
        const prj = project.get(p._id);
        if (prj) {
            p.online = true;
            p.closed = prj.closed;
            p.error = prj.error;
            p.values = _.cloneDeep(prj.variables);
        }
        p.device.forEach(d=>{
            //放进来会让返回包体积太大
            //d.device = p.device.find(v=>v._id.equals(d.device_id));
            //d.element = p.element.find(v=>v._id.equals(d.device.element_id));

            const dvc = device.get(d._id)
            if (dvc) {
                d.online = true;
                d.closed = dvc.closed;
                d.error = dvc.error;
                d.values = _.cloneDeep(dvc.variables);
            }
        })
        //delete p.device;
        //delete p.element;
    },
});