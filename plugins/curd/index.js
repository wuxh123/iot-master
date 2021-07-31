const mongo = require_plugin("mongodb");
const _ = require("lodash");

/**
 * 创建接口
 * @param col
 * @param {function(*)|undefined} before
 * @returns {(function(*): Promise<void>)}
 */
exports.create = function (col, before) {
    return async ctx => {
        if (typeof before === 'function')
            before(ctx)

        const body = ctx.request.body;
        //body.user_id = ctx.state.user._id; 由前端来补充
        const ret = await mongo.db.collection(col).insertOne(body);
        ctx.body = {data: ret.insertedId};

        //后续执行
        process.nextTick(async () => {
            //记录创建事件
            await mongo.db.collection(col + '_event').insertOne({
                [col + '_id']: ctx.params._id,
                event: '创建',
                user_id: ctx.state.user && ctx.state.user._id,
            });
        });
    }
}

exports.setting = function (col, before) {
    return async ctx => {
        if (typeof before === 'function')
            before(ctx)

        const body = ctx.request.body;
        delete body._id; //ID不能修改，MongoDB会报错
        const ret = await mongo.db.collection(col).findOneAndUpdate({_id: ctx.params._id}, {$set: body});
        ctx.body = {data: ret};

        //后续执行
        process.nextTick(async () => {
            //复制原始数据
            const value = _.clone(ret.value);
            delete value._id;
            value[col + '_id'] = ctx.params._id;

            //备份原始的数据
            await mongo.db.collection(col + '_history').insertOne(value)

            //计算差异
            const modify = {};
            for (let k in body) {
                if (!body.hasOwnProperty(k)) continue;
                if (!_.isEqual(body[k], ret.value[k]))
                    modify[k] = body[k];
            }

            //记录修改事件
            await mongo.db.collection(col + '_event').insertOne({
                [col + '_id']: ctx.params._id,
                event: '修改',
                data: modify,
                user_id: ctx.state.user && ctx.state.user._id,
            })
        });
    }
}

exports.detail = function (col, before) {
    return async ctx => {
        if (typeof before === 'function')
            before(ctx)

        const res = await mongo.db.collection(col).findOne({_id: ctx.params._id});
        ctx.body = {data: res}
    }
}

exports.delete = function (col, before) {
    return async ctx => {
        if (typeof before === 'function')
            before(ctx)

        const res = await mongo.db.collection(col).findOneAndDelete({_id: ctx.params._id});
        ctx.body = {data: res}

        //后续执行
        process.nextTick(async () => {
            const value = _.clone(res.value);
            delete value._id;
            value[col + '_id'] = ctx.params._id;

            //备份删除的数据
            await mongo.db.collection(col + '_deleted').insertOne(value);
            //记录删除事件
            await mongo.db.collection(col + '_event').insertOne({
                [col + '_id']: ctx.params._id,
                event: '删除',
                user_id: ctx.state.user && ctx.state.user._id,
            });
        });
    }
}


exports.list = function (col, keywords, before) {
    return async ctx => {
        if (typeof before === 'function')
            before(ctx)

        const body = ctx.request.body;

        //条件
        let $match = {};
        if (body.keyword) {
            const $or = [];
            for (let k in keywords) {
                if (!keywords.hasOwnProperty(k)) return;

                let value = body.keyword;
                if (keywords[k] === 'string') {
                    value = {$regex: body.keyword}
                }
                if (keywords[k] === 'number') {
                    value = parseInt(body.keyword)
                }
                $or.push({[k]: value});
            }

            if ($or.length > 1) {
                $match.$or = $or;//[{name: {$regex: query.keyword}}, {address: query.keyword}, {port: query.keyword}]
            } else if ($or.length === 1) {
                //长度为1，则不需要or了
                Object.assign($match, $or[0]);
            }
        }

        if (body.filter) {
            body.filter.forEach(f => {
                if (Array.isArray(f.value)) {
                    if (f.value.length > 0)
                        $match[f.key] = f.value.length === 1 ? f.value[0] : {$in: f.value};
                } else {
                    $match[f.key] = f.value;
                }
            })
        }

        const pipeline = [{$match}];
        const stages = [
            {$match},
            {$count: 'total'},
            {
                $lookup: {
                    from: col,
                    as: 'data',
                    pipeline
                }
            }
        ]

        //排序
        if (body.sort) {
            const sort = body.sort.filter(s => s.value);
            if (sort.length) {
                const $sort = {};
                sort.forEach(s => {
                    $sort[s.key] = s.value === 'ascend' ? 1 : -1;
                });
                pipeline.push({$sort});
            }
        }

        //分页
        if (body.pageIndex && body.pageSize) {
            pipeline.push({$skip: (body.pageIndex - 1) * body.pageSize});
            pipeline.push({$limit: body.pageSize});
        }

        //查询
        const res = await mongo.db.collection(col).aggregate(stages).toArray();
        if (res.length > 0) {
            ctx.body = res[0];
        } else {
            ctx.body = {total: 0, data: []}
        }
    }
}