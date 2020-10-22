package web

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

func peer(writer http.ResponseWriter, request *http.Request)  {
	key := request.URL.Query()["key"]
	log.Println(key)

	//TODO 获取链接
	//link, err := core.GetLink()


	websocket.Handler(func(ws *websocket.Conn) {
		//peer := core.NewPeer(ws, nil)
		//peer.Receive()
	}).ServeHTTP(writer, request)
	//ctx.Abort()
}

