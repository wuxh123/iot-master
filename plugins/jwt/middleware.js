const jwt = require("koa-jwt");
const index = require('./index');

module.exports = function () {
    return jwt(index.jwtMiddlewareConfig());
}