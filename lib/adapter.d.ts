import {Tunnel} from "./tunnel";

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
     * @param {number} code
     * @param {number} address
     * @param {number} length
     * @returns {Promise<ArrayBuffer>}
     */
    read(slave, code, address, length)

    /**
     * 写单个线圈或寄存器
     * @param {number} slave
     * @param {number} code
     * @param {number} address
     * @param {boolean|number|[]} value
     * @returns {Promise<>}
     */
    write(slave, code, address, value)


    /**
     * 处理数据
     * @param {Buffer|string} data
     */
    handle(data)
}


export declare function create(tunnel: Tunnel, options: Object);