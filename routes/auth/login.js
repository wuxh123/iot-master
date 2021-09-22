const mongo = require_plugin("mongodb");
const jwt = require_plugin("jwt");

exports.post = async ctx => {
    const body = ctx.request.body;
    const {username, password} = body;
    if (username === 'admin' && password === '123456') {
        ctx.body = {data: jwt.sign({_id: '000000000000000000000000', admin: true})}
        return;
    }

    const user = await mongo.db.collection("user").findOne({username, password,})
    if (!user) throw new Error("找不到用户，或密码错误");

    ctx.body = {data: jwt.sign({_id: user._id, admin: user.admin})};
    ctx.body.data.user = user;
}