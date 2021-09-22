const curd = require_plugin("mongodb/curd");
exports.delete = exports.get = curd.delete("alarm");