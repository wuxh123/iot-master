package api

import (
	"git.zgwit.com/iot/dtu-admin/dtu"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func mqtt(ctx *gin.Context)  {
	websocket.Handler(func(ws *websocket.Conn) {
		//设置二进制模式
		ws.PayloadType = websocket.BinaryFrame
		dtu.Hive().Receive(ws)
	}).ServeHTTP(ctx.Writer, ctx.Request)
	//ctx.Abort()
}
