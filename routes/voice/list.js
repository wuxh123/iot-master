const curd = require_plugin("mongodb/curd");
exports.post = curd.list("voice", {
    joins: [{
        from: 'alarm',
        fields: ['name', 'content']
    }, {
        from: 'company',
        fields: ['name']
    }]
});