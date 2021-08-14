const nodeDiskInfo = require("node-disk-info");
const os = require("os");

exports.get = async ctx => {
    const info = await nodeDiskInfo.getDiskInfo()
    ctx.body = {
        data: info.map(disk => {
            if (os.platform() !== "win32") {
                disk.blocks *= 1024;
                disk.used *= 1024;
                disk.available *= 1024;
            }
            return {
                name: disk.filesystem,
                size: disk.blocks,
                used: disk.used,
                free: disk.available,
                usage: parseFloat(disk.capacity),
                mounted: disk.mounted,
            }
        })
    }
}