const os = require("os");

function getTimes(){
    let idle = 0;
    let total = 0;
    os.cpus().forEach(cpu=>{
        for (let type in cpu.times) {
            total += cpu.times[type]
        }
        idle += cpu.times.idle
    });
    return {total, idle};
}

function getUsage(interval){
    return new Promise(((resolve, reject) => {
        let start = getTimes();
        setTimeout(function () {
            let end = getTimes();
            let idle = end.idle - start.idle
            let total = end.total - start.total
            let usage = (10000 - Math.round(10000 * idle / total)) / 100

            return resolve(usage)
        }, interval)
    }));
}

exports.get = async ctx => {
    const interval = parseInt(ctx.query.interval || 1000)
    ctx.body = {
        data: {
            usage: await getUsage(interval),
            count: os.cpus().length,
            model: os.cpus()[0].model,
            speed: os.cpus()[0].speed,
        }
    }
}