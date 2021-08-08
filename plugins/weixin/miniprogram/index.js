const got = require("got");
const modules = require("./modules");

exports.auth = require("./auth");
exports.pay = require("./pay");

//遍历modules.js，获得URL
Object.keys(modules).forEach((m, k) => {
  const api = (exports[k] = {});
  Object.keys(m).forEach((v, k) => {
    api[k] = async function(access_token, params) {
      //TODO 先获取AccessToken
      return await got.post(`${config.domain}/${v}`, {
        searchParams: {
          access_token
        },
        body, 
	params,
        json: true
      });
    };
  });
});
