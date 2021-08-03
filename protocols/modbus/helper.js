/**
 * PLC组态地址，需要去掉万位，并减1
 * 00001 ~ 10000 ：离散量输出，线圈 DO，BO，读01，写05，15
 * 10001 ~ 20000 ：离散量输入，触点 DI，BI，读02
 * 20001 ~ 30000 ：浮点寄存器，IEEE754 读03，04，写06，16
 * 40001 ~ 50000 ：保持寄存器 AO，RH，读03，写06，16
 * 30001 ~ 40000 ：输入寄存器 AI，RI，读04
 * 50001 ~ 60000 ：ASCII字符 读03，04，写06，16
 */

exports.parseReadAddress = function (address) {
    const type = address.substr(0, 2)
    let code = 0x03;
    switch (type) {
        case 'DO':
        case 'BO':
            code = 0x01; // 05 15
            break;
        case 'DI':
        case 'BI':
            code = 0x02;
            break;
        case 'AO':
        case 'RH':
            code = 0x03; // 06 16
            break;
        case 'AI':
        case 'RI':
            code = 0x04;
            break;
        default:
            throw new Error("不支持的地址" + address)
    }
    return {
        code: code,
        address: parseInt(address.substring(2))
    }
}

exports.parseWriteAddress = function (address) {
    const type = address.substr(0, 2)
    let code = 0x03;
    switch (type) {
        case 'DO':
        case 'BO':
            code = 0x05; // 15
            break;
        case 'AO':
        case 'RH':
            code = 0x06; // 16
            break;
        default:
            throw new Error("不支持的地址" + address)
    }
    return {
        code: code,
        address: parseInt(address.substring(2))
    }
}

exports.convertWriteCode = function (code, multi) {
    let c = 0;
    switch (code) {
        case 1:
            c = 5;
            break;
        case 3:
            c = 6;
            break;
        default:
            throw new Error("该功能码不支持写入");
    }
    if (multi)
        c += 10;
    return c;
}

/**
 * 布尔数组压缩成二进制
 * @param {boolean[]|Uint8Array} data
 * @returns {Buffer}
 */
exports.booleanArrayToBuffer = function (data) {
    const size = parseInt((data.length - 1) / 8 + 1);
    const buf = Buffer.alloc(1 + size);
    buf[0] = size; //字节数
    for (let i = 0; i < data.length; i++) {
        buf[parseInt(i / 8) + 1] |= data[i] ? 0x80 >> (i % 8) : 0;
    }
    return buf;
}

/**
 * 数组转成Modbus数据
 * @param data
 * @returns {Buffer}
 */
exports.arrayToBuffer = function (data) {
    let buf;

    if (Array.isArray(data)) {
        const typ = typeof data[0];
        if (typ === 'boolean') {
            buf = exports.booleanArrayToBuffer(data);
        } else {
            //默认字类型：WORD
            const size = data.length * 2;
            buf = Buffer.allocUnsafe(1 + size);
            buf[0] = size; //字节数
            for (let i = 0; i < data.length; i++) {
                buf.writeUInt16BE(data[i], i * 2 + 1);
            }
        }
    } else if (data instanceof Buffer) {
        const size = data.length;
        buf = Buffer.allocUnsafe(1 + size);
        buf[0] = size; //字节数
        data.copy(buf, 1);
    } else if (data instanceof Uint8Array) {
        const size = data.length;
        buf = Buffer.allocUnsafe(1 + size);
        buf[0] = size; //字节数
        for (let i = 0; i < data.length; i++) {
            buf.writeUInt8(data[i], i + 1);
        }
    } else if (data instanceof Uint16Array) {
        const size = data.length * 2;
        buf = Buffer.allocUnsafe(1 + size);
        buf[0] = size; //字节数
        for (let i = 0; i < data.length; i++) {
            buf.writeUInt16BE(data[i], i * 2 + 1);
        }
    } else if (data instanceof Uint32Array) {
        const size = data.length * 4;
        buf = Buffer.allocUnsafe(1 + size);
        buf[0] = size; //字节数
        for (let i = 0; i < data.length; i++) {
            buf.writeUInt32BE(data[i], i * 4 + 1);
        }
    } else if (data instanceof Float32Array) {
        const size = data.length * 4;
        buf = Buffer.allocUnsafe(1 + size);
        buf[0] = size; //字节数
        for (let i = 0; i < data.length; i++) {
            buf.writeFloatBE(data[i], i * 4 + 1);
        }
    } else if (data instanceof Float64Array) {
        const size = data.length * 8;
        buf = Buffer.allocUnsafe(1 + size);
        buf[0] = size; //字节数
        for (let i = 0; i < data.length; i++) {
            buf.writeDoubleBE(data[i], i * 8 + 1);
        }
    } else {
        const size = data.length * 2;
        buf = Buffer.allocUnsafe(1 + size);
        buf[0] = size; //字节数
        for (let i = 0; i < data.length; i++) {
            buf.writeUInt16BE(data[i], i * 2 + 1);
        }
    }
    return buf;
}

exports.crc16 = function (buf) {
    let crc = 0xFFFF;
    let odd;

    for (let i = 0; i < buf.length; i++) {
        crc = crc ^ buf[i];

        for (let j = 0; j < 8; j++) {
            odd = crc & 0x0001;
            crc = crc >> 1;
            if (odd) {
                crc = crc ^ 0xA001;
            }
        }
    }

    return crc;
};

exports.lrc = function lrc(buf) {
    let lrc = 0;
    for (let i = 0; i < buf.length; i++) {
        lrc += buf[i] & 0xFF;
    }

    return ((lrc ^ 0xFF) + 1) & 0xFF;
};
