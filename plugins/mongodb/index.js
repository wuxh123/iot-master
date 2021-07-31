const mongodb = require('mongodb');
const MongoClient = mongodb.MongoClient;

const defaultOptions = {
    host: 'localhost',
    port: 27017,
    db: 'test',
    authSource: 'admin',
    max: 100,
    min: 1,
    maxIdleTimeMS: 60000
};


exports.config = function (opts) {
    const cfg = Object.assign({}, defaultOptions, opts);
    let mongoUrl;
    if (cfg.username && cfg.password) {
        mongoUrl = `mongodb://${cfg.username}:${cfg.password}@${cfg.host}:${cfg.port}/${cfg.db}?authSource=${cfg.authSource}`
    } else {
        mongoUrl = `mongodb://${cfg.host}:${cfg.port}/${cfg.db}`
    }
    const options = {
        useNewUrlParser: true,
        useUnifiedTopology: true,
        maxPoolSize: cfg.max,
        minPoolSize: cfg.min,
        maxIdleTimeMS: cfg.maxIdleTimeMS,
    };

    //TODO 处理错误
    (async function () {
        await MongoClient.connect(mongoUrl, options, function (err, client) {
            if (err) {
                console.log(err);
                return
            }
            exports.client = client;
            exports.db = client.db(cfg.db);
        });
    })();
}

/**
 * 转ObjectId
 * @type {ObjectId}
 */
exports.ObjectId = mongodb.ObjectId;

/**
 * ID转日期
 * @param _id
 * @returns {Date}
 * @constructor
 */
exports.ObjectIdToDate = function (_id) {
    return new Date(parseInt(_id.toString().substring(0, 8), 16) * 1000);
};

/**
 * 日期转ObjectId
 * @param {number} time 时间
 * @returns {ObjectId}
 * @constructor
 */
exports.DateToObjectId = function (time) {
    return mongodb.ObjectId((~~(+time / 1000)).toString(16) + '0'.repeat(16));
};

/**
 * 选择表
 * @param col
 * @returns {mongodb.Collection<Document>}
 */
exports.collection = function (col) {
    return exports.db.collection(col);
};

exports.bodyToStages = function (body) {
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
    return stages;
}