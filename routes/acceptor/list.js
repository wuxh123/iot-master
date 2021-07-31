const curd = require_plugin("curd");
exports.post = curd.list("acceptor", {
    name: 'string',
    type: 'string',
    address: 'string',
    port: 'number',
});