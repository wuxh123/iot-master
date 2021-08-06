const curd = require_plugin("curd");
exports.post = curd.list("subscribe", {
    join: {
        from: 'user',
        fields: ['name']
    }
});