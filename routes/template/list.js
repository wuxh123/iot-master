const curd = require_plugin("curd");
exports.post = curd.list("template", {
    name: 'string',
    type: 'string',
    address: 'string',
    port: 'number',
});