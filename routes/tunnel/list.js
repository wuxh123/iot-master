const tunnel = require("../../lib/tunnel");

const curd = require_plugin("mongodb/curd");
exports.post = curd.list("tunnel", {
    joins: [{
        from: 'acceptor',
        fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }]
});