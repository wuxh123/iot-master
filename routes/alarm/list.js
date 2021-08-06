const curd = require_plugin("curd");
exports.post = curd.list("device", {
    join:{
        from: 'device',
        local: 'device_id',
        fields: ['name']
    }
});