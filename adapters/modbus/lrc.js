module.exports = function lrc(buf) {
    let lrc = 0;
    for (let i = 0; i < buf.length; i++) {
        lrc += buf[i] & 0xFF;
    }

    return ((lrc ^ 0xFF) + 1) & 0xFF;
};
