module.exports = function crc16(buf) {
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
