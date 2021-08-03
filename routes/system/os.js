const os = require("os");

exports.get = async ctx => {
    ctx.body = {
        data: {
            //os: await osu.os.oos(),
            oss: os.version(),
            platform: os.platform(),
            uptime: os.uptime(),
            //ip: osu.os.ip(),
            hostname: os.hostname(),
            type: os.type(),
            arch: os.arch(),
        }
    }
}