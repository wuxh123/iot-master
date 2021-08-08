const got = require("got");
const config = require("../config")();

function sendGet(url, params) {
    return got(url, {searchParams: params});
}

module.exports = {

    getAuthorizeUrl(redirect_uri, state, scope='base'){
        const uri = encodeURIComponent(redirect_uri);
        const st = encodeURIComponent(state);
        return `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${config.official.appid}&redirect_uri=${uri}&response_type=code&scope=snsapi_${scope}&state=${st}#wechat_redirect`
    },

    async getAccessToken(code) {
        const { appid, secret } = config.official;
        return await sendGet(`${config.domain}/sns/oauth2/access_token`, {
            appid,
            secret,
            code,
            grant_type: "authorization_code"
        });
    },

    async refreshAccessToken(refresh_token) {
        const { appid } = config.official;
        return await sendGet(`${config.domain}/sns/oauth2/refresh_token`, {
            appid,
            grant_type: "refresh_token",
            refresh_token
        });
    },

    async checkAccessToken(access_token, openid) {
        const { appid } = config.official;
        return await sendGet(`${config.domain}/sns/auth`, {
            access_token, openid
        });
    },

    async getUserInfo(access_token, openid, lang = 'zh_CN') {
        return await sendGet(`${config.domain}/sns/userinfo`, {
            access_token, openid, lang
        });
    },

}