const mongo = require_plugin("mongodb");
const _ = require("lodash");

/**
 * 创建接口
 * @param col
 * @param {object|undefined} options
 * @returns {(function(*): Promise<void>)}
 */
exports.create = function (col, options) {
    options = options || {};

    return async ctx => {
        if (options.before)
            await options.before(ctx)

        const body = ctx.request.body;
        delete body._id;
        //body.user_id = ctx.state.user._id; 由前端来补充
        const ret = await mongo.db.collection(col).insertOne(body);
        ctx.body = {data: ret.insertedId};

        //后续执行
        process.nextTick(async () => {
            //记录创建事件
            await mongo.db.collection('event').insertOne({
                target: col,
                [col + '_id']: ctx.params._id,
                event: '创建',
                user_id: ctx.state.user && ctx.state.user._id,
            });
        });
    }
}

exports.setting = function (col, options) {
    options = options || {};

    return async ctx => {
        if (options.before)
            await options.before(ctx)

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
            await mongo.db.collection('event').insertOne({
                target: col,
                [col + '_id']: ctx.params._id,
                event: '修改',
                data: modify,
                user_id: ctx.state.user && ctx.state.user._id,
            })
        });
    }
}

exports.detail = function (col, options) {
    options = options || {};

    return async ctx => {
        if (options.before)
            await options.before(ctx)

        const res = await mongo.db.collection(col).findOne({_id: ctx.params._id});
        if (res) ctx.body = {data: res}
        else ctx.body = {error: '找不到数据'}
    }
}

exports.delete = function (col, options) {
    options = options || {};

    return async ctx => {
        if (options.before)
            await options.before(ctx)

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
            await mongo.db.collection('event').insertOne({
                target: col,
                [col + '_id']: ctx.params._id,
                event: '删除',
                user_id: ctx.state.user && ctx.state.user._id,
            });
        });
    }
}


exports.list = function (col, options) {
    options = options || {};

    return async ctx => {
        if (options.before)
            await options.before(ctx)

        const body = ctx.request.body || {};

        let pipeline = [
            {$match: body.filter || {}},
            {$sort: body.sort || {_id: -1}},
            {$skip: body.skip || 0},
            {$limit: body.limit || 20},
            //TODO project
        ];

        if (options.pipeline) {
            pipeline = pipeline.concat(pipeline)
        }

        function addJoin(join) {
            const local = join.local || join.from + '_id';
            const foreign = join.foreign || '_id';
            const $lookup = {
                from: join.from,
                as: join.as || join.from,
                let: {id: '$'+local},
                pipeline: [
                    {$match: {$expr: {$eq: ["$"+foreign, "$$id"]}}},
                ],
            };
            if (join.fields && join.fields.length) {
                const $project = {};
                join.fields.forEach(f=>$project[f]=1);
                $lookup.pipeline.push({$project});
            }

            pipeline.push({$lookup});
            pipeline.push({$unwind: {path: '$' + $lookup.as, preserveNullAndEmptyArrays: true}});

            if (join.replace) {
                pipeline.push({$addFields: {[$lookup.as + '.' + col + '_id']: '$_id'}})
                pipeline.push({$replaceRoot: {newRoot: '$' + $lookup.as}});
            }
        }

        options.join && addJoin(options.join)
        options.joins && options.joins.forEach(addJoin)

        //支持参数中的fields
        const fields = body.fields || options.fields;
        if (fields && fields.length) {
            const $project = {};
            fields.forEach(f=>$project[f]=1);
            pipeline.push({$project});
        }

        const stages = [
            {$match: body.filter || {}},
            {$count: 'total'}, //计算总数
            {
                //数据查询
                $lookup: {
                    from: col,
                    as: 'data',
                    pipeline
                }
            }
        ]

        //查询
        const res = await mongo.db.collection(col).aggregate(stages).toArray();
        ctx.body = res.length ? res[0] : {total: 0, data: []}
    }
}