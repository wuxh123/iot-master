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


module.exports = function (opts) {
    const cfg = Object.assign({}, defaultOptions, mongoConfig, opts);
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
            exports.client = client;
            exports.db = client.db(cfg.db);
        });
    })();
}


exports.ObjectId = mongodb.ObjectId;

exports.ObjectIdToDate = function (_id) {
    return new Date(parseInt(_id.toString().substring(0, 8), 16) * 1000);
};

exports.DateToObjectId = function (time) {
    return mongodb.ObjectId((~~(+time / 1000)).toString(16) + '0'.repeat(16));
};

exports.collection = function (col) {
    return exports.db.collection(col);
};