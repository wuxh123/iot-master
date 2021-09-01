const mongo = require_plugin("mongodb")
const vms = require_plugin("tencentcloud/vms");

exports.notice = async function (alarm) {
    log.info(alarm, 'notice');

    const date = new Date();
    const time = date.getHours() * 60 + date.getMinutes();

    //记录数据库
    const ret = await mongo.db.collection("alarm").insertOne(alarm)
    alarm._id = ret.insertedId;


    const filter = {
        enable: true,
        names: alarm.name,
        //level: {$lte: alarm.level},
        $or: [
            {$expr: {$lt: ['$start', '$end']}, start: {$lte: time}, end: {$gte: time},},
            //加入隔天判断
            {$expr: {$gt: ['$start', '$end']}, $or: [{start: {$lte: time}}, {end: {$gte: time}}],}
        ],
    }

    //union查询订阅，先项目，然后是分组，最后是公司
    const pipeline = []
    pipeline.push({$match: Object.assign({}, {project_id: alarm.project_id}, filter)});

    //查询分组订阅
    if (alarm.group_id)
        pipeline.push({
            $unionWith: {
                coll: 'subscribe',
                pipeline: [{$match: Object.assign({}, {group_id: alarm.group_id}, filter)}]
            }
        })

    //查询企业订阅
    if (alarm.company_id)
        pipeline.push({
            $unionWith: {
                coll: 'subscribe',
                pipeline: [{$match: Object.assign({}, {company_id: alarm.company_id}, filter)}]
            }
        });

    //查询用户
    pipeline.push({
        $lookup: {
            from: 'user',
            as: 'user',
            localField: 'user_id',
            foreignField: '_id'
        }
    })
    pipeline.push({$unwind: {path: '$user'}}) //, preserveNullAndEmptyArrays: true} 找不到用户 就过滤掉 或者 TODO 删除无效订阅
    //pipeline.push({$replaceRoot: {newRoot: '$user'}})

    //查出所有订阅
    const subs = await mongo.db.collection("subscribe").aggregate(pipeline).toArray();

    //获取联系方式并去重
    let smsSubs = subs.filter(sub => sub.sms && sub.user.cellphone).map(sub => sub.user.cellphone);
    smsSubs = [...new Set(smsSubs)];
    log.info(smsSubs, '短信通知')

    let emailSubs = subs.filter(sub => sub.email && sub.user.email).map(sub => sub.user.email);
    emailSubs = [...new Set(emailSubs)];
    log.info(emailSubs, '邮件通知');

    let wxSubs = subs.filter(sub => sub.weixin && sub.user.wx && sub.user.wx.official).map(sub => sub.user.wx.official.openid);
    wxSubs = [...new Set(wxSubs)];
    log.info(wxSubs, '微信通知');

    let voiceSubs = subs.filter(sub => sub.voice && sub.user.cellphone).map(sub => sub.user.cellphone);
    voiceSubs = [...new Set(voiceSubs)];
    log.info(voiceSubs, '语音通知');

    //记录电话列表，电话回调中处理
    if (voiceSubs.length) {
        alarm.cellphone = voiceSubs;
        await mongo.db.collection("alarm").updateOne({_id: ret.insertedId}, {$set: {cellphone: voiceSubs}})
        //先向第一个人通知
        await sendVoice(alarm, voiceSubs[0]);
    }
}

async function sendVoice(alarm, cellphone) {
    log.info({cellphones: alarm.cellphone, cellphone}, '语音通知');
    const res = await mongo.db.collection("voice").insertOne({alarm_id: alarm._id, company_id:alarm.company_id, cellphone: cellphone,})
    try {
        const resp = await vms.send([alarm.project_name || '', alarm.device_name || '', alarm.content], cellphone, 2, "")
        await mongo.db.collection("voice").updateOne({_id: res.insertedId}, {$set: {callid: resp.SendStatus.CallId}})
    } catch (e) {
        await mongo.db.collection("voice").updateOne({_id: res.insertedId}, {$set: {error: e.message}})
        //通知下一个
        const index = alarm.cellphone.indexOf(cellphone);
        if (index < alarm.cellphone.length - 1)
            await sendVoice(alarm, alarm.cellphone[index + 1]);
    }
}

exports.handleVoiceNotify = async function (body) {
    /*
    成功：
           "result": "0",
           "accept_time": "1470197211",
           "call_from": "",
           "callid": "xxxxxx",
           "end_calltime": "1470197221",
           "fee": "1",
           "mobile": "13xxxxxxxxx",
           "nationcode": "86",
           "start_calltime": "1470197196"
    失败：
            callid: '0bd403fe-ffd6-11eb-935a-5254002ab735',
            failure_code: 11,
            failure_reason: '其他',
            call_from: '089831335739',
            mobile: '18621890022',
            nationcode: '86'
    */

    await mongo.db.collection("voice").updateOne({callid: body.callid}, {$set: body})
    if (body.result === "0") {
        //通知成功了，就结束
        return;
    }

    const voice = await mongo.db.collection("voice").findOne({callid: body.callid});
    if (!voice) return;
    const alarm = await mongo.db.collection("alarm").findOne({_id: voice.alarm_id});
    if (!alarm) return;
    if (alarm.closed) return;//已经关闭的提示不需要再打电话
    const cnt = await mongo.db.collection("voice").countDocuments({
        alarm_id: voice.alarm_id,
        cellphone: voice.cellphone
    })
    //0正常，1未接听，2异常
    if (cnt < 3 && body.result === "1") {
        //未接听，5分钟后重拨
        setTimeout(() => {
            sendVoice(alarm, voice.cellphone).then(() => {
            }).catch(log.error);
        }, 5 * 60 * 1000); //5分钟

    } else {
        //通知下一个
        const index = alarm.cellphone.indexOf(voice.cellphone);
        if (index < alarm.cellphone.length - 1)
            await sendVoice(alarm, alarm.cellphone[index + 1]);
        //else 没有可以通知了
    }
}