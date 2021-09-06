const {protocols} = require("../../lib/plugin");

exports.get = async ctx => {
    ctx.body = {data: protocols}
}