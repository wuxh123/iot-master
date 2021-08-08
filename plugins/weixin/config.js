let config = {
    miniprogram: {
        appid: "",
        secret: "",
        mch_id: "",
        paykey: "",
        refundCAPath: "",
    },
    official: {
        appid: "",
        secret: "",
        mch_id: "",
        paykey: "",
        refundCAPath: "",
    },
    domain: "https://api.weixin.qq.com",
    paydomain: "https://api.mch.weixin.qq.com",
    spbill_create_ip: ""
};

module.exports = function (c) {
    if (c)
        Object.assign(config, c);
    return config;
}
