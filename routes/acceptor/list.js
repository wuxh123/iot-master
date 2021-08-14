const acceptor = require("../../lib/acceptor");

const curd = require_plugin("curd");
exports.post = curd.list("acceptor", {
    after: async ctx => {
        ctx.body.data.forEach(a => {
            const acc = acceptor.getAcceptor(a._id);
            if (acc) {
                a.online = true;
                a.closed = acc.closed;
                a.error = acc.error;
            }
        })
    }
});