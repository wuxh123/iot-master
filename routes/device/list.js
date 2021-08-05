const curd = require_plugin("curd");
exports.post = curd.list("device", {
    joins: [{
        from: 'element',
        fields: ['name']
    }, {
        from: 'tunnel',
        fields: ['name', 'sn']
    }]
});