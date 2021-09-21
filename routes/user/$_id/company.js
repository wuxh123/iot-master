const mongo = require_plugin("mongodb");
exports.get = (async ctx => {
    const user_id = ctx.params._id;
    const ret = await mongo.db.collection("company").aggregate([
        //自己的
        {$match: {user_id: user_id}},
        //成员的
        {
            $unionWith: {
                coll: 'member', pipeline: [
                    {$match: {user_id: user_id}},
                    {
                        $lookup: {
                            from: "company",
                            let: {id: "$company_id"},
                            pipeline: [
                                {$match: {$expr: {$eq: ["$_id", "$$id"]}}},
                            ],
                            as: "company"
                        }
                    },
                    {$unwind: {path: "$company"}},
                    {$addFields: {"company.member_id": "$_id", "company.admin": "$admin"}},
                    {$replaceRoot: {newRoot: "$company"}},
                ]
            }
        },
    ]).toArray();

    ctx.body = {data: ret};
});