const path = require('path');
const acceptor = require("./acceptor");

const mongodb = require_plugin("mongodb");

//数据目录
global.data_path = process.env.DATA_PATH || path.join(process.cwd(), 'data');

//加载配置
require('./config');

//启用插件功能
require('./plugin');

//启用网页
require('./web')

// 恢复 接收器
const acceptors = mongodb.collection("acceptor").find({}).toArray();
acceptors.forEach(options => {
    const ac = acceptor.create(options);


});


