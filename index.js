const path = require('path');

//数据目录
global.data_path = process.env.DATA_PATH || path.join(process.cwd(), 'data');

//加载配置
require('./lib/setting');

//启用插件功能
require('./lib/plugin');


//TODO 1.连接数据库

//TODO 2.恢复acceptor（连接成功后）

//TODO 3.启动Web服务

//TODO 4.启动透传服务
