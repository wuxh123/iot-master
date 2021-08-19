const jwt = require_plugin('jwt');

exports.get = (async ctx => {

    ctx.body = {data: jwt.sign({_id: ctx.params._id, proxy: true})};
});