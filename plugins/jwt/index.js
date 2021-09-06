const jwt = require('jsonwebtoken');
const koaJwt = require('koa-jwt');
const _ = require("lodash");

const defaultOptions = {
    secret: 'jwt-secret',
    expiresIn: 24 * 60 * 60, //s 一天
    getToken(ctx, opts) {
        return ctx.query[opts.queryToken || 'token'];
    }
};

const cfg = load_config("jwt")
const options = _.defaultsDeep({}, cfg, defaultOptions)

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

