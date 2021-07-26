declare class Adapter {
    /**
     * @type Tunnel
     */
    tunnel;

    options;

    /**
     * 构造函数
     * @param {Tunnel} tunnel
     * @param {Object} options
     */
    constructor(tunnel, options)

    /**
     * 读取数据
     * @param {number} slave
     * @param {string} address
     * @param {number} length
     * @returns {Promise<Uint16Array|Uint8Array>}
     */
    read(slave, address, length)

    /**
     * 写单个线圈或寄存器
     * @param {number} slave
     * @param {string} address
     * @param {boolean|number} value
     * @returns {Promise<>}
     */
    write(slave, address, value)


    /**
     * 写入多个线圈或寄存器
     * @param {number} slave
     * @param {string} address
     * @param {boolean[]|number[]} data
     * @returns {Promise<>}
     */
    writeMany(slave, address, data)

    /**
     * 处理数据
     * @param {Buffer|string} data
     */
    handle(data)
}