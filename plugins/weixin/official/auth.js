const got = require("got");
const config = require("../config")();

function sendGet(url, params) {
  return got(url, {searchParams: params});
}

module.exports = {
  /**
   * 获取全局AccessToken
   * 获取小程序全局唯一后台接口调用凭据（access_token）。调调用绝大多数后台接口时都需使用 access_token，开发者需要进行妥善保存。
   * 默认2小时
   * @param {String} grant_type 默认值：client_credential
   * @param {Boolean} cache 是否使用自带缓存.默认:true
   * @example getAccessToken({grant_type,cache})
   */
  async getAccessToken(grant_type = "client_credential") {
    const { appid, secret } = config.official;
    return await sendGet(`${config.domain}/cgi-bin/token`, {
        appid,
        secret,
        grant_type
      });
  },

};
