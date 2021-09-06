const mongodb = require('mongodb');
const _ = require("lodash");
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

let cfg = load_config("mongodb");
_.defaultsDeep(cfg, defaultOptions);


/**
 *
 * @type {MongoClient}
 */
exports.client = null

/**
 *
 * @type {mongodb.Db}
 */
exports.db = null;


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
    MongoClient.connect(mongoUrl, options, function (err, client) {
        if (err) {
            console.log(err);
            return
        }
        exports.client = client;
        exports.db = client.db(cfg.db);

        //执行等待
        callbacks.forEach(c => c());
    });

const callbacks = [];

exports.ready = function (callback) {
    if (this.client) callback();
    else callbacks.push(callback);
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
 * @returns {mongodb.Collection}
 */
exports.collection = function (col) {
    if (!exports.db) {
        throw new Error("MongoDB数据库未连接")
    }
    return exports.db.collection(col);
};
