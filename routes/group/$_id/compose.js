const curd = require_plugin("curd");
exports.get = curd.compose("group", {
    joins: [{
        from: 'user',
        fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }]
});