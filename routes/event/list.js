const curd = require_plugin("mongodb/curd");
exports.post = curd.list("event", {
    join: {
        from: 'user',
        local: 'user_id',
        fields: ['name', 'avatar']
    }
});