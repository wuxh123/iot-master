const got = require("got");
const config = require("../config")();

function sendGet(url, params) {
  return got(url, {searchParams: params}).json();
}

module.exports = {
  /**
   * 微信服务端登录
   * @param {String} js_code
   * @param {String} grant_type 默认值:authorization_code
   * @example code2Session({js_code,grant_type})
   */
  async code2Session(js_code, grant_type = "authorization_code") {
    const { appid, secret } = config.miniprogram;
    return await sendGet(`${config.domain}/sns/jscode2session`, {
        appid,
        secret,
        js_code,
        grant_type
      });
  },

  /**
   * 获取全局AccessToken
   * 获取小程序全局唯一后台接口调用凭据（access_token）。调调用绝大多数后台接口时都需使用 access_token，开发者需要进行妥善保存。
   * 默认2小时
   * @param {String} grant_type 默认值：client_credential
   * @param {Boolean} cache 是否使用自带缓存.默认:true
   * @example getAccessToken({grant_type,cache})
   */
  async getAccessToken(grant_type = "client_credential") {
    const { appid, secret } = config.miniprogram;
    return await sendGet(`${config.domain}/cgi-bin/token`, {
        appid,
        secret,
        grant_type
      });
  },

  /**
   * 获取该用户的 UnionId
   * @example getPaidUnionId(option)
   */
  async getPaidUnionId(option) {
    return await sendGet(`${config.domain}/wxa/getpaidunionid`, option);
  }
};
