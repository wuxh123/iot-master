package web

import (
	"github.com/kataras/iris/v12"
	"golang.org/x/net/websocket"
	"log"
)

func peer(ctx iris.Context)  {
	key := ctx.URLParam("key")
	log.Println(key)

	//TODO 获取链接
	//link, err := core.GetLink()


	websocket.Handler(func(ws *websocket.Conn) {
		//peer := core.NewPeer(ws, nil)
		//peer.Receive()
	}).ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	//ctx.Abort()
}

