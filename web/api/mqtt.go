package api

import (
	"git.zgwit.com/iot/dtu-admin/dbus"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func monitor(ctx *gin.Context)  {
	websocket.Handler(func(conn *websocket.Conn) {
		//转入MQTT
		dbus.Hive().Receive(conn)
	}).ServeHTTP(ctx.Writer, ctx.Request)
}
