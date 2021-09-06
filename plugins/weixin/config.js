const _ = require("lodash");
let defaultConfig = {
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

const cfg = load_config("weixin");

let config = _.defaultsDeep({}, cfg, defaultConfig);


module.exports = function () {
    return config;
}
