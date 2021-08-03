const nodeDiskInfo = require("node-disk-info");

exports.get = async ctx => {
    const info = await nodeDiskInfo.getDiskInfo()
    ctx.body = {
        data: info.map(disk => {
            return {
                name: disk.filesystem,
                size: disk.blocks,
                used: disk.used,
                free: disk.available,
                usage: disk.capacity,
                mounted: disk.mounted,
            }
        })
    }
}