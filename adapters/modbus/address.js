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