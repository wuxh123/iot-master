const mongo = require_plugin("mongodb");
const curd = require_plugin("curd");
exports.post = curd.create("member", {
    before: async ctx => {
        const body = ctx.request.body;
        const m = await mongo.db.collection("member").findOne({user_id:body.user_id, company_id: body.company_id});
        if (m) throw new Error("已经存在");
    }
});