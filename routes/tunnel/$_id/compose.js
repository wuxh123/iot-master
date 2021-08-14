const acceptor = require("../../../lib/acceptor");

const curd = require_plugin("curd");
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
        const prj = acceptor.getTunnel(p._id);
        if (prj) {
            p.online = true;
            p.closed = prj.closed;
            p.error = prj.error;
        }
    },
});