package web

import (
	"git.zgwit.com/zgwit/dtu-admin/core"
	"golang.org/x/net/websocket"
	"net/http"
)

func mqtt(writer http.ResponseWriter, request *http.Request)  {
	websocket.Handler(func(ws *websocket.Conn) {
		//设置二进制模式
		ws.PayloadType = websocket.BinaryFrame
		core.Hive().Receive(ws)
	}).ServeHTTP(writer, request)
	//ctx.Abort()
}
