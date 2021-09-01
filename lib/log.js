const pino = require('pino');
const path = require("path");
const _ = require("lodash");

const defaultOptions = {
    level: 'info',
    base: {}, //清空 pid 和 hostname
    prettyPrint: {
        levelFirst: true,
        translateTime: 'yyyy-mm-dd HH:MM:ss',
        singleLine: true,
    }
};

const filename =  path.join(global.data_path, "log.config");
const options = _.defaultsDeep({}, require(filename), defaultOptions);

//module.exports = pino(options);
global.log = pino(options);