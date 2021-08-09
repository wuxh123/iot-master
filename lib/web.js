const path = require("path");
const Koa = require("koa");
const koaStatic = require('koa-static');
const koaLogger = require('koa-logger');
const koaRouter = require('koa-router');
const koaCors = require("@koa/cors");
const koaBodyParser = require('koa-bodyparser');
const koaEasyWS = require('koa-easy-ws');


const app = new Koa();

//启用日志
app.use(koaLogger())

//静态文件在前，每次都会先检查文件是否存在
//app.use(koaStatic('www'));

//跨域问题
app.use(koaCors());

app.use(koaEasyWS());

app.use(koaBodyParser({enableTypes: ['json', 'form', 'xml']}));

app.use(async (ctx, next) => {
    //ctx.state.user = {_id: mongo.ObjectId("000000000000000000000000")}; //测试用
    try {
        await next();
    } catch (err) {
        console.error(err); //打印错误，供调试
        ctx.body = {error: err.message};
    }
});


//启用JWT
const jwtMiddleware = require_plugin("jwt").middleware();

const r = koaRouter();
r.use(jwtMiddleware.unless({path: [/^\/doc/, /^\/api\/auth/, /^\/open\/auth/, /\/watch/,]}))
r.use(require("./routes").middleware());
r.use(require("./open").middleware());

app.use(r.routes());
app.use(r.allowedMethods());

//静态文件
app.use(koaStatic(path.join(__dirname, '..', 'www')));

app.listen(8008);