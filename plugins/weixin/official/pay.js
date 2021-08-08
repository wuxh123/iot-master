const got = require("got");
const config = require("../config")();
const fs = require("fs");
const {
  xmlToJson,
  jsonToXml,
  signToMD5,
  sortByASCII,
  getNonceStr
} = require("../utils");

/**
 * 内部请求代理对象
 * @param {*} url 请求地址
 * @param {*} params 支付参数
 * @param {*} requestOption request对象参数
 */
async function sendPost(url, params, option = {}) {
  const { paykey } = config.official;
  const payParams = sortByASCII(params);
  const sign = signToMD5(payParams, paykey);
  payParams.sign = sign;

  const payParamsXml = jsonToXml(payParams).trim();

  const payResult = await got(
    Object.assign(
      {
        uri: url,
        method: "POST",
        searchParams: params,
        body: payParamsXml,
        contentType: "text/xml",
        convertJSON: false,
        json: false
      },
      option
    )
  );

  return xmlToJson(payResult);
}

module.exports = {
  /**
   * 统一下单
   * @param {*} option 支付参数对象
   */
  async unifiedOrder(option) {
    const { appid, mch_id } = config.official;

    const payParams = Object.assign(option, {
      appid,
      mch_id,
      nonce_str: option.nonce_str || getNonceStr(),
      trade_type: option.trade_type || "JSAPI",
      //notify_url
    });

    return sendPost(`${config.paydomain}/pay/unifiedorder`, payParams);
  },

  /**
   * 查询订单
   * @param {*} option 查询参数对象
   */
  async orderQuery(option) {
    const { appid, mch_id } = config.official;
    const payParams = Object.assign(option, {
      appid,
      mch_id,
      nonce_str: option.nonce_str || getNonceStr()
    });

    return sendPost(`${config.paydomain}/pay/orderquery`, payParams);
  },

  /**
   * 关闭订单
   * @param {*} option 关闭参数对象
   */
  async closeOrder(option) {
    const { appid, mch_id } = config.official;
    const payParams = Object.assign(option, {
      appid,
      mch_id,
      nonce_str: option.nonce_str || getNonceStr()
    });

    return sendPost(`${config.paydomain}/pay/closeorder`, payParams);
  },

  /**
   * 申请退款
   * @param {*} option 退款参数对象
   */
  async refund(option) {
    const { appid, mch_id, refundCAPath } = config.official;
    const payParams = Object.assign(option, {
      appid,
      mch_id,
      nonce_str: option.nonce_str || getNonceStr()
    });

    return sendPost(`${config.paydomain}/secapi/pay/refund`, payParams, {
      agentOptions: {
        pfx: fs.readFileSync(refundCAPath),
        passphrase: mch_id
      }
    });
  },

  /**
   * 查询退款
   * @param {*} option 退款查询对象
   */
  async refundQuery(option) {
    const { appid, mch_id } = config.official;
    const payParams = Object.assign(option, {
      appid,
      mch_id,
      nonce_str: option.nonce_str || getNonceStr()
    });

    return sendPost(`${config.paydomain}/pay/refundquery`, payParams);
  },

  /**
   * 下载对账单
   * @param {*} option 下载参数对象
   */
  async downloadBill(option) {
    const { appid, mch_id } = config.official;
    const payParams = Object.assign(option, {
      appid,
      mch_id,
      nonce_str: option.nonce_str || getNonceStr()
    });

    return sendPost(`${config.paydomain}/pay/downloadbill`, payParams);
  }
};
