const curd = require_plugin("curd");
exports.post = curd.list("tunnel", {
    name: 'string',
    type: 'string',
    address: 'string',
    port: 'number',
});