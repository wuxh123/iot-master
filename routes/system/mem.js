const os = require("os");

exports.get = async ctx => {
    ctx.body = {
        data: {
            size: os.totalmem(),
            free: os.freemem(),
            used: os.totalmem() - os.freemem(),
            //usage: (os.totalmem() - os.freemem()) * 100 / os.totalmem(),
            usage: (10000 - Math.round(10000 * os.freemem() / os.totalmem())) / 100
        }
    }
}