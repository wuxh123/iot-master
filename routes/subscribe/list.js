const curd = require_plugin("mongodb/curd");
exports.post = curd.list("subscribe", {
    joins: [{
        from: 'user',
        fields: ['name', 'avatar']
    }, {
        from: 'project',
        fields: ['name'],
    }, {
        from: 'group',
        fields: ['name'],
    }, {
        from: 'company',
        fields: ['name'],
    }]
});