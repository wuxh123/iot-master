require("./lib/init");

//避免异常错误导致程序退出
process.on('uncaughtException', function (err) {
    //打印出错误
    console.error(err);
    //打印出错误的调用栈方便调试
    //console.log(err.stack);
});
