const curd = require_plugin("mongodb/curd");
exports.post = curd.list("company", {
    join: {
        from: 'user',
        fields: ['name']
    }
});