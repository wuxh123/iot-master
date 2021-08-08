const crypto = require("crypto");
const cryptoRandomString = require("crypto-random-string");
const convert = require("xml-js");

/**
 * 生成32位随机字符串
 */
function getNonceStr() {
  return cryptoRandomString({ length: 32 });
}

/**
 * json转xml
 * @param {*} jsonObj  需要转换的json对象
 */
function jsonToXml(jsonObj) {
  const result = convert.js2xml(jsonObj, {
    compact: true,
    ignoreComment: false,
    spaces: 4
  });
  return `<xml>${result}</xml>`;
}

/**
 * xml转json
 * @param {*} xmlObj 需要转换的xml对象
 */
function xmlToJson(xmlObj) {
  var result = convert.xml2js(xmlObj, { compact: false, spaces: 4 });
  let jsonResult = {};
  result.elements[0].elements.forEach(item => {
    if (item.elements)
      jsonResult[item.name] =
        item.elements[0].text !== undefined
          ? item.elements[0].text
          : item.elements[0].cdata;
  });
  return jsonResult;
}
/**
 * 根据ASCII排序指定对象
 * @param {*} param 排序的对象
 */
function sortByASCII(param) {
  const keys = Object.keys(param).sort();
  const sortObj = {};
  keys.forEach(key => {
    sortObj[key] = param[key];
  });
  return sortObj;
}
/**
 * 生成微信支付签名（MD5格式）
 * @param {*} params 签名对象
 * @param {*} payKey 商户密钥
 * @example {appid,mch_id,nonce_str}
 */
function signToMD5(params, payKey) {
  const tempParams = JSON.parse(JSON.stringify(sortByASCII(params)));

  if (tempParams.sign) {
    delete tempParams.sign;
  }

  let tempStr = "";
  const paramKeys = Object.keys(tempParams);
  paramKeys.forEach(key => {
    tempStr = `${tempStr}&${key}=${params[key]}`;
  });
  tempStr = tempStr.substring(1) + "&key=" + payKey;

  const md5 = crypto.createHash("md5");
  md5.update(tempStr);

  return md5.digest("hex").toUpperCase();
}

function verifyPaySign() {}

/**
 * 校验微信小程序用户信息
 * @param {String} session_key
 * @param {String} rawData
 * @param {String} signature
 * @example checkSign({session_key,rawData,signature})
 */
function checkUserSign(option) {
  const { session_key, rawData, signature } = option;
  if (
    session_key === undefined ||
    rawData === undefined ||
    signature === undefined ||
    rawData === undefined
  ) {
    return false;
  }
  const serverSign = crypto
    .createHash("sha1")
    .update(rawData + session_key)
    .digest("hex");
  return signature === serverSign;
}

/**
 * 解析微信小程序encryptedData
 * @param {*} encryptedData 微信小程序的userInfo.encryptedData
 * @param {String} sessionKey wx.login登录后的sessionKey
 * @param {String} iv 加密算法的初始向量
 * @example decryptData({encryptedData,sessionKey,iv})
 */
function decryptData(option) {
  const { encryptedData, sessionKey, iv, appid } = option;
  return wxBizDataCrypt(encryptedData, iv, sessionKey, appid);
}

function wxBizDataCrypt(encryptedData, iv, sessionKey, appid) {
  // base64 decode
  sessionKey = Buffer.from(sessionKey, "base64");
  encryptedData = Buffer.from(encryptedData, "base64");
  iv = Buffer.from(iv, "base64");

  try {
    // 解密
    var decipher = crypto.createDecipheriv("aes-128-cbc", sessionKey, iv);
    // 设置自动 padding 为 true，删除填充补位
    decipher.setAutoPadding(true);
    var decoded = decipher.update(encryptedData, "binary", "utf8");
    decoded += decipher.final("utf8");

    decoded = JSON.parse(decoded);
  } catch (err) {
    throw new Error("Illegal Buffer");
  }

  if (decoded.watermark.appid !== appid) {
    throw new Error("Illegal Buffer");
  }

  return decoded;
}

module.exports = {
  getNonceStr,
  jsonToXml,
  xmlToJson,
  sortByASCII,
  signToMD5,
  checkUserSign,
  decryptData
};
