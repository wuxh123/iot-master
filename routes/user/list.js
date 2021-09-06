const curd = require_plugin("mongodb/curd");
exports.post = curd.list("user");