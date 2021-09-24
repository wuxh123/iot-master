const mongo = require_plugin("mongodb");
const jwt = require_plugin("jwt");

const crypto = require("crypto");

function md5(str){
    return crypto.createHash("md5").update(str).digest('hex');
}

exports.post = async ctx => {
    const body = ctx.request.body;
    let {username, password} = body;

    let user = await mongo.db.collection("user").findOne({username})
    if (!user) {
        if (username === 'admin') {
            //创建默认管理员
            const ret = await mongo.db.collection("user").insertOne({
                name: '超级管理员',
                admin: true,
                username: username,
                password: md5(md5('123456')),
                enable: true
            });
            user = await mongo.db.collection("user").findOne({_id: ret.insertedId});
        } else {
            throw new Error("找不到用户");
        }
    }
    if (!user.enable)
        throw new Error("用户被禁用");

    //二次加密
    password = md5(password);
    if (password !== user.password)
        throw new Error("密码错误");

    ctx.body = {data: jwt.sign({_id: user._id, admin: user.admin})};
    ctx.body.data.user = user;
}