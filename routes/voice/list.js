const curd = require_plugin("curd");
exports.post = curd.list("voice", {
    joins: [{
        from: 'alarm',
        fields: ['name', 'content']
    }, {
        from: 'company',
        fields: ['name']
    }]
});