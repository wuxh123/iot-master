const path = require('path');

//数据目录
global.data_path = process.env.DATA_PATH || path.join(process.cwd(), 'data');

//加载配置
require('./config');

//启用插件功能
require('./plugin');

//启用协议
require('./protocol');

//启用网页
if (config.web.enable)
    require('./web')

//启动服务
const acceptor = require('./acceptor');

//恢复 接收器
const mongo = require_plugin("mongodb");
mongo.ready(function (){
    mongo.db.collection("acceptor").find({enable: true}).toArray().then(res => {
        res.forEach(model => {
            try {
                acceptor.create(model);
            } catch (err) {
                console.error(err.message)
            }
        });
    });
})



