const curd = require_plugin("curd");
exports.post = curd.list("tunnel", {
    joins: [{
        from: 'acceptor',
        fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }]
});