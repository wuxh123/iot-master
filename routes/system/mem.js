const osu = require("node-os-utils");

exports.get = async ctx => {
    ctx.body = {data: await osu.mem.info()}
}