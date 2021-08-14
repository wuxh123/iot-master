const mongodb = require("mongodb");

function promote(data) {
    if (Array.isArray(data)) {
        data.forEach((v, k) => {
            data[k] = promote(v);
        });
    } else if (typeof data === "object" && data) {
        Object.keys(data).forEach(k => {
            const d = data[k];
            // ObjectId
            if (/_id$/.test(k)) {
                if (Array.isArray(d)) {
                    d.forEach((v, k) => {
                        if (typeof v === "string") {
                            try {
                                d[k] = mongodb.ObjectId(v);
                            } catch (e) {}
                        }
                    });
                } else if (typeof d === "string") {
                    try {
                        data[k] = mongodb.ObjectId(d);
                    } catch (e) {}
                }
            } else if (k === "created" || k === "updated" || /_date$/.test(k)) {
                if (typeof d === "string") {
                    try {
                        data[k] = new Date(new Date(d).toISOString());
                    } catch (e) {}
                } else if (typeof d === "object" && d) {
                    Object.keys(d).forEach(k => {
                        if (typeof d[k] === "string") {
                            try {
                                data[k] = new Date(new Date(d).toISOString());
                            } catch (e) {}
                        }
                    });
                }
            } else {
                data[k] = promote(data[k]);
            }
        });
    }
    return data;
}

function demote(data) {
    if (Array.isArray(data)) {
        data.forEach((v, k) => {
            data[k] = demote(v);
        });
    } else if (typeof data === "object" && data) {
        Object.keys(data).forEach(k => {
            const d = data[k];
            // ObjectId
            if (/_id$/.test(k)) {
                if (Array.isArray(d)) {
                    d.forEach((v, k) => {
                        if (v instanceof mongodb.ObjectId) {
                            d[k] = v.toHexString();
                        }
                    });
                } else if (d instanceof mongodb.ObjectId) {
                    data[k] = d.toHexString();
                }
            } else if (k === "created" || k === "updated" || /_date$/.test(k)) {
                if (d instanceof Date) {
                    data[k] = d.toISOString();
                } else if (typeof d === "object" && d) {
                    Object.keys(d).forEach(k => {
                        if (d[k] instanceof Date) {
                            data[k] = d.toISOString();
                        }
                    });
                }
            } else {
                data[k] = demote(data[k]);
            }
        });
    }
    return data;
}

module.exports = function(options = {}) {
    return async function(ctx, next) {
        if (typeof ctx.request.query.query === "string") {
            try {
                ctx.request.query.query = JSON.parse(ctx.request.query.query);
            } catch (e) {}
        }
        if (typeof ctx.request.query.sort === "string") {
            try {
                ctx.request.query.sort = JSON.parse(ctx.request.query.sort);
            } catch (e) {}
        }

        //类型升级，方便处理
        promote(ctx.request.query);

        if (ctx.method === "POST" || ctx.method === "PUT") {
            promote(ctx.request.body);
        }

        //路由参数
        promote(ctx.params);

        //状态（JWT payload）
        promote(ctx.state.user);

        await next();

        //类型降级，接下来交给joi-router检验output
        //demote(ctx.body);
    };
};
