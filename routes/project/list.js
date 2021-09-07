const project = require("../../lib/project");

const curd = require_plugin("mongodb/curd");
exports.post = curd.list("project", {
    joins: [{
        from: 'template',
        fields: ['name']
    }, {
        from: 'company',
        fields: ['name']
    }, {
        from: 'group',
        fields: ['name']
    }, {
        from: 'user',
        fields: ['name']
    }]
});