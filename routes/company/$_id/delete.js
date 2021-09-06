const curd = require_plugin("mongodb/curd");
exports.delete = exports.get = curd.delete("company");

//TODO after钩子，删除相关项目，解绑相关通道