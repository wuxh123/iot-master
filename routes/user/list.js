const curd = require_plugin("curd");
exports.post = curd.list("user", {
    name: 'string',
});