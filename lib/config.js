const fs = require("fs");
const path = require("path");
const _ = require("lodash");

const defaultOptions = {
    web: {
        enable: true,
        port: 8888,
        cors: true,
        logger: false,
    },
    traverser: {
        enable: true,
        port: 1843,
    },
    debug: false,
};


const filename =  path.join(global.data_path, "config");
global.config = _.defaultsDeep({}, require(filename), defaultOptions);