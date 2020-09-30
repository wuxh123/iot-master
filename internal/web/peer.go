package web

import (
	"git.zgwit.com/zgwit/iot-admin/internal/core"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"log"
)

func peer(ctx *gin.Context)  {
	key := ctx.Query("key")
	log.Println(key)

	//TODO 获取链接
	//link, err := core.GetLink()


	websocket.Handler(func(ws *websocket.Conn) {
		peer := core.NewPeer(ws, nil)
		peer.Receive()
	}).ServeHTTP(ctx.Writer, ctx.Request)
	//ctx.Abort()
}

