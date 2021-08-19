const tencentcloud = require("tencentcloud-sdk-nodejs");
// 导入 VMS 模块的 client models
const vmsClient = tencentcloud.vms.v20200902.Client;
/* 实例化要请求 VMS 的 client 对象 */
let client;// = new vmsClient();
let options = {};

exports.config = function (opts) {
    options = opts;
    client = new vmsClient(opts);
}

exports.send = function (set, cellphone, times, sess) {
    console.log('打电话', set, cellphone)
    return client.SendTtsVoice({
        // 模板 ID，必须填写在控制台审核通过的模板 ID，可登录 [语音消息控制台] 查看模板 ID
        TemplateId: options.TemplateId, //"1043772",
        TemplateParamSet: set, //["7652","hello","world"],
        /* 被叫手机号码，采用 e.164 标准，格式为+[国家或地区码][用户号码]
         * 例如：+8613711112222，其中前面有一个+号，86为国家码，13711112222为手机号
         */
        CalledNumber: "+86"+cellphone, //"+8615161515197",
        /* 在语音控制台添加应用后生成的实际SdkAppid，示例如1400006666 */
        VoiceSdkAppid: options.VoiceSdkAppid, //"1400550599",
        /* 播放次数，可选，最多3次，默认2次 */
        PlayTimes: times || 3,
        /* 用户的 session 内容，腾讯 server 回包中会原样返回 */
        SessionContext: sess,
    })
}
