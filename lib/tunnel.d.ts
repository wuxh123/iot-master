import {EventEmitter} from "events";

declare class Tunnel extends EventEmitter {
    /**
     * 连接
     * @type Socket
     */
    conn;

    /**
     * 过滤器
     * @type {*[]}
     */
    filters;

    /**
     * 适配器
     */
    adapter;

    /**
     * 透传通道
     * @type Socket
     */
    traversal;

    options;

    /**
     * 初始化
     * @param {Socket} conn 网络连接
     * @param {Object} options 参数
     */
    constructor(conn, options)

    /**
     * 开启透传
     * socket关闭，即自动关闭
     * 不支持断线重连
     * @param {Socket} socket
     */
    traverse(socket)

    /**
     * 发送数据
     * @param {Buffer|string} data
     * @returns boolean
     */
    write(data)

    /**
     * 关闭通道
     */
    close()
}