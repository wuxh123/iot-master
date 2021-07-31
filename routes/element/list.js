const curd = require_plugin("curd");
exports.post = curd.list("element", {
    name: 'string',
    type: 'string',
    address: 'string',
    port: 'number',
});