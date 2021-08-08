const jwt = require('jsonwebtoken');
const koaJwt = require('koa-jwt');

const defaultOptions = {
    secret: 'jwt-secret',
    expiresIn: 24 * 60 * 60 //s 一天
};

let options = Object.assign({}, defaultOptions);

//const jwtConfig = require('../../../../jwt.config');

exports.config = function (opts) {
    options = Object.assign({}, defaultOptions, opts);
}

exports.jwtMiddlewareConfig = function () {
    return options;
}

exports.sign = function (obj) {
    const token = jwt.sign(obj, options.secret, {
        expiresIn: options.expiresIn
    })

    return {
        token,
        expiresAt: Math.floor(Date.now() / 1000) + options.expiresIn
    }
}


exports.middleware = function () {
    return koaJwt(options);
}

