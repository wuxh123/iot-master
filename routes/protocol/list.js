const {protocols} = require("../../lib/protocol");

exports.get = async ctx => {
    ctx.body = {data: protocols}
}