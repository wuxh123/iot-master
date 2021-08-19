const {handleVoiceNotify} = require("../../lib/notice");

exports.post = async ctx=>{
    const result = ctx.request.body.voiceprompt_callback; // || ctx.request.body.voice_failure_callback;
    //console.log("voice_callback", result);
    if (result)
        await handleVoiceNotify(result)
    ctx.body = {
        "result": 0,
        "errmsg": "OK"
    }
}