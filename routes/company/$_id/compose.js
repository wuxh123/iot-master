const curd = require_plugin("mongodb/curd");
exports.get = curd.compose("company", {
    join: {
        from: 'user',
        fields: ['name']
    }
});