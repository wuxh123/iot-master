const path = require("path");
const Koa = require("koa");
const koaStatic = require('koa-static');
const send = require('koa-send');
const logger = require('koa-logger');
const Router = require('koa-router');
const bodyParser = require('koa-bodyparser');
const unless = require('koa-unless');
const ws = require('koa-easy-ws');

const app = new Koa();

//启用日志
if (config.web.logger)
    app.use(logger())

//静态文件在前，每次都会先检查文件是否存在
//app.use(koaStatic('www'));

//跨域问题
if (config.web.cors) {
    const cors = require("@koa/cors");
    app.use(cors());
}

app.use(ws());

app.use(bodyParser({enableTypes: ['json', 'form', 'xml']}));

app.use(async (ctx, next) => {
    //ctx.state.user = {_id: mongo.ObjectId("000000000000000000000000")}; //测试用
    try {
        await next();
    } catch (err) {
        if (config.debug)
            log.error(err.message); //打印错误，供调试
        ctx.status = err.status || 200;
        ctx.body = {error: err.message};
    }
});


//启用JWT
const jwtMiddleware = require_plugin("jwt").middleware();

const r = Router();
r.use(jwtMiddleware.unless({path: [/^\/doc/, /^\/api\/auth/, /^\/api\/voice\/callback/, /^\/open\/auth/]}))
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
const defaultPage = async ctx => await send(ctx, 'index.html', {root: path.join(__dirname, '..', 'www')});
defaultPage.unless = unless;
app.use(defaultPage.unless({path: [/^\/api\//, /^\/open\//]}));

//开始监听web
app.listen(config.web.port);
if (config.debug)
    log.info(config.web, 'listen');