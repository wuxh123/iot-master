const curd = require_plugin("curd");
exports.post = curd.list("company", {
    join: {
        from: 'user',
        fields: ['name']
    }
});