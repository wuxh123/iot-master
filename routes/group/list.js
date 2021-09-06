const curd = require_plugin("mongodb/curd");
exports.post = curd.list("group", {
    joins: [{
        from: 'user',
        fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }]
});