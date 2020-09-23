package api

import (
	"git.zgwit.com/zgwit/iot-admin/internal/channel"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func mqtt(ctx *gin.Context)  {
	websocket.Handler(func(ws *websocket.Conn) {
		//设置二进制模式
		ws.PayloadType = websocket.BinaryFrame
		channel.Hive().Receive(ws)
	}).ServeHTTP(ctx.Writer, ctx.Request)
	//ctx.Abort()
}
