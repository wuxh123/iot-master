package web

import (
	"git.zgwit.com/zgwit/iot-admin/core"
	"github.com/kataras/iris/v12"
	"golang.org/x/net/websocket"
)

func mqtt(ctx iris.Context)  {
	websocket.Handler(func(ws *websocket.Conn) {
		//设置二进制模式
		ws.PayloadType = websocket.BinaryFrame
		core.Hive().Receive(ws)
	}).ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	//ctx.Abort()
}
