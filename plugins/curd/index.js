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

        if (options.after)
            await options.after(ctx)

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
        let update = body;
        if (!update.$set && !update.$unset && !update.$push)
            update = {$set: update};

        const ret = await mongo.db.collection(col).findOneAndUpdate({_id: ctx.params._id}, update);
        const obj = await mongo.db.collection(col).findOne({_id: ctx.params._id});
        ctx.body = {data: obj};

        if (options.after)
            await options.after(ctx)

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
        if (!res)
            throw new Error("找不到数据")
        ctx.body = {data: res}

        if (options.after)
            await options.after(ctx)
    }
}

exports.delete = function (col, options) {
    options = options || {};

    return async ctx => {
        if (options.before)
            await options.before(ctx)

        const res = await mongo.db.collection(col).findOneAndDelete({_id: ctx.params._id});
        ctx.body = {data: res}

        if (options.after)
            await options.after(ctx)

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


exports.compose = function (col, options) {
    options = options || {};

    return async ctx => {
        if (options.before)
            await options.before(ctx)

        const body = ctx.request.body || {};

        /**
         *
         * @type {[]}
         */
        let pipeline = [
            {$match: {_id: ctx.params._id}},
        ];

        function addJoin(join) {
            const local = join.local || join.from + '_id';
            const foreign = join.foreign || '_id';
            const as = join.as || join.from;
            if (join.fields && join.fields.length) {
                const $lookup = {
                    from: join.from,
                    as: as,
                    let: {id: '$' + local},
                    pipeline: [
                        {$match: join.filter || {$expr: {$eq: ["$" + foreign, "$$id"]}}},
                    ],
                };
                const $project = {};
                join.fields.forEach(f => $project[f] = 1);
                $lookup.pipeline.push({$project});

                pipeline.push({$lookup});
            } else {
                //简单$lookup 支持对象数组作为条件的查询
                const $lookup = {
                    from: join.from,
                    as: as,
                    localField: local,
                    foreignField: foreign,
                };
                pipeline.push({$lookup});
            }
            if (!join.noUnwind)
                pipeline.push({$unwind: {path: '$' + as, preserveNullAndEmptyArrays: true}});
        }

        options.join && addJoin(options.join)
        options.joins && options.joins.forEach(addJoin)

        //支持参数中的fields
        const fields = body.fields || options.fields;
        if (fields && fields.length) {
            const $project = {};
            fields.forEach(f => $project[f] = 1);
            pipeline.push({$project});
        }

        //查询
        const res = await mongo.db.collection(col).aggregate(pipeline).toArray();
        if (!res.length)
            throw new Error("找不到记录")
        ctx.body = {data: res[0]}

        if (options.after)
            await options.after(ctx)
    }
}


exports.list = function (col, options) {
    options = options || {};

    return async ctx => {
        if (options.before)
            await options.before(ctx)

        const body = ctx.request.body || {};

        let pipeline = [
            ...(ctx.state.stages || []),
            {$match: body.filter || {}},
            {$sort: body.sort || {_id: -1}},
            {$skip: body.skip || 0},
            {$limit: body.limit || 20},
        ];

        if (options.pipeline) {
            pipeline = pipeline.concat(pipeline)
        }

        function addJoin(join) {
            const local = join.local || join.from + '_id';
            const foreign = join.foreign || '_id';
            const as = join.as || join.from;
            if (join.fields && join.fields.length) {
                const $lookup = {
                    from: join.from,
                    as: as,
                    let: {id: '$' + local},
                    pipeline: [
                        {$match: join.filter || {$expr: {$eq: ["$" + foreign, "$$id"]}}},
                    ],
                };
                const $project = {};
                join.fields.forEach(f => $project[f] = 1);
                $lookup.pipeline.push({$project});

                pipeline.push({$lookup});
            } else {
                //简单$lookup 支持对象数组作为条件的查询
                const $lookup = {
                    from: join.from,
                    as: as,
                    localField: local,
                    foreignField: foreign,
                };
                pipeline.push({$lookup});
            }
            if (!join.noUnwind)
                pipeline.push({$unwind: {path: '$' + as, preserveNullAndEmptyArrays: true}});
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
            fields.forEach(f => $project[f] = 1);
            pipeline.push({$project});
        }

        const stages = [
            ...(ctx.state.stages || []),
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

        if (options.after)
            await options.after(ctx)
    }
}