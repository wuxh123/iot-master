const curd = require_plugin("mongodb/curd");
exports.post = curd.list("alarm", {
    joins: [{
        from: 'device',
        fields: ['name']
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