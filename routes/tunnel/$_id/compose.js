const tunnel = require("../../../lib/tunnel");

const curd = require_plugin("mongodb/curd");
exports.get = curd.compose("tunnel", {
    joins: [{
        from: 'acceptor',
        fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }],
    after: async ctx => {
        const p = ctx.body.data;
        const tnl = tunnel.get(p._id);
        if (tnl) {
            p.online = tnl.online;
        }
    },
});