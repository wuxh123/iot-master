const internals = {};

exports.check = function(internal, callback) {
    if (!internals.hasOwnProperty(internal)) {
        internals[internal] = {
            handler: setInterval(function () {
                const now = Date.now();
                internals[internal].callbacks.forEach(callback=> {
                    //加入异常，避免被中断
                    try {
                        callback(now)
                    } catch (e) {
                        //log.error(e)
                    }
                });
            }, internal),
            internal,
            callbacks: []
        };
    }
    internals[internal].callbacks.push(callback);

    return {
        cancel(){
            const index = internals[internal].callbacks.findIndex(cb=>cb===callback);
            if (index > -1) {
                internals[internal].callbacks.splice(index, 1);
                if (internals[internal].callbacks.length === 0) {
                    clearInterval(internals[internal].handler)
                    delete internals[internal];
                }
            }
        }
    }
}
