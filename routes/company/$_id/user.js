const curd = require_plugin("curd");
exports.post = curd.list("member", {
    before: ctx=>{
        ctx.state.stages = [
            {$match: {company_id: ctx.params._id}},
            {
                $lookup: {
                    from: 'user',
                    localField: 'user_id',
                    foreignField: '_id',
                    as: 'user'
                }
            },
            {$unwind: {path: '$user'}},
            {$addFields: {'user.member_id': '$_id', 'user.admin': '$admin'}},
            {$replaceRoot: {newRoot: '$user'}},
        ];
    }
});