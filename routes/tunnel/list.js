const acceptor = require("../../lib/acceptor");

const curd = require_plugin("curd");
exports.post = curd.list("tunnel", {
    after: async ctx => {
        ctx.body.data.forEach(p => {
            const prj = acceptor.getTunnel(p._id);
            if (prj) {
                p.online = true;
                //p.closed = prj.closed;
                //p.error = prj.error;
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