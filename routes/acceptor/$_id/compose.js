const acceptor = require("../../../lib/acceptor");

const curd = require_plugin("mongodb/curd");
exports.get = curd.compose("acceptor", {
    after: async ctx => {
        const a = ctx.body.data;
        const acc = acceptor.getAcceptor(a._id);
        if (acc) {
            a.online = true;
            a.closed = acc.closed;
            a.error = acc.error;
        }
    }
});