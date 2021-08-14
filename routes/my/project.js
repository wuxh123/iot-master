const mongo = require_plugin("mongodb");
const project = require('../../lib/project')
const device = require('../../lib/device')

exports.get = (async ctx => {
    const user_id = ctx.state.user._id;

    //一、先找到公司的项目（自己的 和 管理员的）
    let results = await mongo.db.collection("company").aggregate([
        //自己的公司
        {$match: {user_id: user_id}},

        //作为成员的公司
        {
            $unionWith: {
                coll: 'member',
                pipeline: [
                    {$match: {user_id: user_id, admin: true}},
                    {
                        $lookup: {
                            from: "company",
                            let: {id: "$company_id"},
                            pipeline: [
                                {$match: {$expr: {$eq: ["$_id", "$$id"]}}}, //只查管理员的
                            ],
                            as: "company"
                        }
                    },
                    {$unwind: {path: "$company"}},
                    {$addFields: {"company.member_id": "$_id"}}, //company.manager: $manager
                    {$replaceRoot: {newRoot: "$company"}},
                ]
            }
        },
        //找到相关项目
        {
            $lookup: {
                from: 'project',
                let: {id: "$_id"},
                pipeline: [
                    {$match: {$expr: {$eq: ["$company_id", "$$id"]}}},
                    {$project: {name: 1, devices: 1, user_id: 1}},
                    {
                        $lookup: {
                            from: 'device',
                            localField: 'devices.device_id',
                            foreignField: '_id',
                            as: 'device',
                        }
                    },
                    {
                        $lookup: {
                            from: 'element',
                            localField: 'device.element_id',
                            foreignField: '_id',
                            as: 'element',
                        }
                    },
                    //补充用户信息
                    {
                        $lookup: {
                            from: "user",
                            let: {id: "$user_id"},
                            pipeline: [
                                {$match: {$expr: {$eq: ["$_id", "$$id"]}}},
                                {$project: {name: 1}}
                            ],
                            as: "user"
                        }
                    },
                    {$unwind: {path: "$user", preserveNullAndEmptyArrays: true}},
                ],
                as: 'projects',
            }
        },
        //标记以上为公司项目
        {$addFields: {"type": "company"}},
    ]).toArray();

    const companies = results.map(c => c._id);
    //二、查找找到作为组长的，且不在上述公司的
    let ret = await mongo.db.collection("group").aggregate([
        {$match: {user_id: user_id, company_id: {$nin: companies}}},
        {
            $lookup: {
                from: 'project',
                let: {id: "$_id"},
                pipeline: [
                    {$match: {$expr: {$eq: ["$group_id", "$$id"]}}},
                    {$project: {name: 1, devices: 1, user_id: 1}},
                    {
                        $lookup: {
                            from: 'device',
                            localField: 'devices.device_id',
                            foreignField: '_id',
                            as: 'device',
                        }
                    },
                    {
                        $lookup: {
                            from: 'element',
                            localField: 'device.element_id',
                            foreignField: '_id',
                            as: 'element',
                        }
                    },
                    //补充用户信息
                    {
                        $lookup: {
                            from: "user",
                            let: {id: "$user_id"},
                            pipeline: [
                                {$match: {$expr: {$eq: ["$_id", "$$id"]}}},
                                {$project: {name: 1}}
                            ],
                            as: "user"
                        }
                    },
                    {$unwind: {path: "$user", preserveNullAndEmptyArrays: true}},
                ],
                as: 'projects',
            }
        },
        //标记以上为分组项目
        {$addFields: {"type": "group"}},
    ]).toArray();

    //拼接数据
    results = results.concat(ret.filter(g => g.projects.length));

    //计算已经查找到的项目
    let projects = results.map(c => c.projects.map(p => p._id));
    let ids = projects.length ? projects.reduce((v1, v2) => v1.concat(v2)) : []; //reduce 空数据会报错

    //三、找到作为负责人的
    ret = await mongo.db.collection("project").aggregate([
        {$match: {user_id: user_id, _id: {$nin: ids}}},
        {$project: {name: 1, devices: 1, user_id: 1}},
        {
            $lookup: {
                from: 'device',
                localField: 'devices.device_id',
                foreignField: '_id',
                as: 'device',
            }
        },
        {
            $lookup: {
                from: 'element',
                localField: 'device.element_id',
                foreignField: '_id',
                as: 'element',
            }
        },
        //补充用户信息
        {
            $lookup: {
                from: "user",
                let: {id: "$user_id"},
                pipeline: [{$match: {$expr: {$eq: ["$_id", "$$id"]}}},],
                as: "user"
            }
        },
        {$unwind: {path: "$user", preserveNullAndEmptyArrays: true}},
    ]).toArray();

    if (ret.length)
        results = results.concat({type: 'project', projects: ret})


    //填充状态
    results.forEach(c => {
        c.projects.forEach(p => {
            //项目状态
            const prj = project.get(p._id)
            if (prj) {
                p.online = true;
                p.closed = prj.closed;
                p.error = prj.error;
                p.values = prj.variables;
            }

            //设备状态
            p.device.forEach(d => {
                const dvc = device.get(d._id);
                if (dvc) {
                    d.online = true;
                    d.closed = dvc.closed;
                    d.error = dvc.error;
                    d.values = dvc.variables;
                }
            })
        })
    })

    ctx.body = {data: results};
});