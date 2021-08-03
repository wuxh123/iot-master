const osu = require("node-os-utils");
const os = require("os");

exports.get = async ctx => {
    ctx.body = {
        data: {
            usage: await osu.cpu.usage(),
            count: os.cpus().length,
            model: os.cpus()[0].model,
            speed: os.cpus()[0].speed,
        }
    }
}