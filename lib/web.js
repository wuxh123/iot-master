const path = require("path");
const Koa = require("koa");
const koaStatic = require('koa-static');
const send = require('koa-send');
const logger = require('koa-logger');
const Router = require('koa-router');
const cors = require("@koa/cors");
const bodyParser = require('koa-bodyparser');
const unless = require('koa-unless');
const ws = require('koa-easy-ws');


const app = new Koa();

//启用日志
app.use(logger())

//静态文件在前，每次都会先检查文件是否存在
//app.use(koaStatic('www'));

//跨域问题
app.use(cors());

app.use(ws());

app.use(bodyParser({enableTypes: ['json', 'form', 'xml']}));

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

const r = Router();
r.use(jwtMiddleware.unless({path: [/^\/doc/, /^\/api\/auth/, /^\/open\/auth/, /\/watch/,]}))
r.use('/api', require("./routes").middleware());
r.use('/open', require("./open").middleware());

app.use(r.routes());
app.use(r.allowedMethods());
//TODO Router中应该添加404判断


//前端静态文件
const staticFiles = koaStatic(path.join(__dirname, '..', 'www'));
staticFiles.unless = unless;
app.use(staticFiles.unless({path: [/^\/api\//, /^\/open\//]}));

//默认发送首页，以支持前端框架的无#路由
app.use(async ctx => {
    await send(ctx, 'index.html', {root: path.join(__dirname, '..', 'www')})
})

//开始监听web
app.listen(8008);