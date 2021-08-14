const jwt = require('jsonwebtoken');
//const jwtConfig = require('../../../../jwt.config');

exports.get = (async ctx => {
    const jwtObj = {
        _id: ctx.params._id,
    };

    const token = jwt.sign(jwtObj, jwtConfig.secret, {
        expiresIn: jwtConfig.expiresIn
    });
    //TODO 失效期改为1天

    ctx.body = {
        data: {
            token,
            jwtObj,
            expiresAt: Math.floor(Date.now() / 1000) + jwtConfig.expiresIn
        }
    };
});