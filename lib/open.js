const fs = require('fs');
const path = require('path');
const pad = require_plugin("mongodb/middleware");
const methods = require("methods");
const send = require('koa-send');

const router = require("koa-joi-router")();
//const KoaJoiRouterDocs = require("koa-joi-router-docs");

function scanRoutes(dir, prefix) {
    const files = fs.readdirSync(dir)
    files.forEach(function (name) {
        const filePath = path.join(dir, name);
        const stat = fs.statSync(filePath);
        if (stat.isFile()) {
            const ext = path.extname(name);
            if (ext === '.js') {
                //TODO 要使用eval或vm，因为文件会更新，模拟加载
                const mod = require(path.resolve(filePath));
                const base = path.basename(name, ext);
                const p = (base === 'index') ? prefix : prefix + base;
                console.log('[open]', p)
                for (const k in mod) {
                    if (mod.hasOwnProperty(k) && methods.indexOf(k) > -1) {
                        if (typeof mod[k] === 'function') {
                            router.route({
                                path: p,
                                method: k,
                                handler: [pad(), mod[k]]
                            })
                        } else {
                            const spec = mod[k];
                            spec.path = p + (spec.path || '');
                            spec.method = k;
                            spec.handler = [pad(), spec.handler]
                            router.route(spec)
                        }
                    }
                }
            } else {
                //发送静态文件
                router.get(prefix + name, async ctx => await send(ctx, name, {root: dir}))
            }
        } else if (stat.isDirectory()) {
            if (name.substr(0, 1) === '$') {
                scanRoutes(filePath, prefix + ':' + name.substr(1) + '/');
            } else {
                scanRoutes(filePath, prefix + name + '/');
            }
        }
    });

    return router;
}

scanRoutes(path.join(__dirname, '..', 'open'), '/');

module.exports = router;

module.exports.reload = function () {
    //清空joi-router和koa-router内容，尚未验证是否生效
    router.routes = [];
    router.router.stack = [];
    scanRoutes(path.join(__dirname, '..', 'open'), '/');
}