const mongo = require_plugin("mongodb");

exports.post = async ctx => {
    const body = ctx.request.body;

    const stages = [];
    let $match = {};
    if (body.keyword) {
        $match.$or = [{name: {$regex: body.keyword}}, {address: body.keyword}, {port: body.keyword}]
    }
    if (body.filter) {
        body.filter.filter(f => f.value.length > 0).forEach(f => {
            $match[f.key] = f.value.length === 1 ? f.value[1] : {$in: f.value};
        })
    }
    stages.push({$match});

    if (body.sort) {
        const $sort = {};
        body.sort.forEach(s => {
            $sort[s.key] = s.value === 'ascend' ? 1 : -1;
        })
        stages.push({$sort})
    }
    if (body.pageIndex && body.pageSize) {
        stages.push({$skip: body.pageIndex * body.pageSize});
        stages.push({$limit: body.pageSize});
    }

    //ctx.request.body
    const res = await mongo.collection("acceptor").aggregate(stages).toArray();
    ctx.body = {data: res}
}