const {protocols} = require("../../lib/adapter");

exports.get = async ctx => {
    ctx.body = {data: protocols}
}