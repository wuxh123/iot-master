const tunnel = require("../../lib/tunnel");

const curd = require_plugin("mongodb/curd");
exports.post = curd.list("tunnel", {
    after: async ctx => {
        ctx.body.data.forEach(p => {
            const tnl = tunnel.get(p._id);
            if (tnl) {
                p.online = tnl.online;
            }
        })
    },
    joins: [{
        from: 'acceptor',
        fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }]
});