const osu = require("node-os-utils");
const os = require("os");

exports.get = async ctx => {
    ctx.body = {
        data: {
            os: await osu.os.oos(),
            platform: os.platform(),
            uptime: os.platform(),
            ip: osu.os.ip(),
            hostname: os.hostname(),
            type: os.type(),
            arch: os.arch(),
        }
    }
}