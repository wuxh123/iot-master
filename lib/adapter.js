class Adapter {
    /**
     * @type Tunnel
     */
    tunnel;

    options = {};

    /**
     * 构造函数
     * @param {Tunnel} tunnel
     * @param {Object} options
     */
    constructor(tunnel, options) {
        this.tunnel = tunnel;
        Object.assign(this.options, options);
    }

    /**
     * 读取数据
     * @param {number} slave
     * @param {string} address
     * @param {number} length
     * @returns {Promise<Uint16Array|Uint8Array>}
     */
    read(slave, address, length) {
    }

    /**
     * 写单个线圈或寄存器
     * @param {number} slave
     * @param {string} address
     * @param {boolean|number} value
     * @returns {Promise<>}
     */
    write(slave, address, value) {
    }


    /**
     * 写入多个线圈或寄存器
     * @param {number} slave
     * @param {string} address
     * @param {boolean[]|number[]} data
     * @returns {Promise<>}
     */
    writeMany(slave, address, data) {
    }

    /**
     * 处理数据
     * @param {Buffer|string} data
     */
    handle(data) {
    }
}

/**
 * 创建适配器
 * @param {Tunnel} tunnel
 * @param {Object} options 协议配置
 * @return {Adapter}
 */
exports.create = function(tunnel, options) {
    //TODO 检查js文件
    return new require('../adapters/' + options.type)(this, options.options);
}