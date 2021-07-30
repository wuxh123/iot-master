const Koa = require("koa");
const logger = require('koa-logger');
const cors = require("@koa/cors");
const bodyParser = require('koa-bodyparser');


const app = new Koa();

//启用日志
app.use(logger())

//跨域问题
app.use(cors());

app.use(bodyParser({enableTypes: ['json', 'form', 'xml']}));

app.use(async (ctx, next) => {
    //ctx.state.user = {_id: mongo.ObjectId("000000000000000000000000")}; //测试用
    try {
        await next();
    } catch (err) {
        console.error(err); //打印错误，供调试
        ctx.body = {err: err.message};
    }
});


//启用JWT
// const jwt = require("koa-jwt");
// app.use(jwt(require("./jwt.config")).unless({path: [/^\/doc/, /^\/wx/, /^\/test/, ]}));

//app.use(routes.middleware());
app.use(require("./routes").middleware());

app.listen(8088);