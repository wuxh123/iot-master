package web

import (
	"git.zgwit.com/zgwit/MyDTU/core"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func mqtt(ctx *gin.Context)  {
	websocket.Handler(func(ws *websocket.Conn) {
		//设置二进制模式
		ws.PayloadType = websocket.BinaryFrame
		core.Hive().Receive(ws)
	}).ServeHTTP(ctx.Writer, ctx.Request)
	//ctx.Abort()
}
