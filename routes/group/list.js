const curd = require_plugin("curd");
exports.post = curd.list("group", {
    name: 'string',
    type: 'string',
    address: 'string',
    port: 'number',
});