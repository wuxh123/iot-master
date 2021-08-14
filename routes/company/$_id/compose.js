const curd = require_plugin("curd");
exports.get = curd.compose("company", {
    join: {
        from: 'user',
        fields: ['name']
    }
});