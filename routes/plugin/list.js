const {plugins} = require("../../lib/plugin");

exports.get = async ctx=>{
    ctx.body = {data: plugins}
}