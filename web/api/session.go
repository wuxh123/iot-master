package api

import (
	"encoding/gob"
	"iot-master/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

var sessionHandler gin.HandlerFunc

func init() {
	//注册 User类型
	gob.Register(&model.User{})

	//初始化Session
	store := memstore.NewStore([]byte("secret"))
	sessionHandler = sessions.Sessions("sess", store)
}

func GetSessionHandler() gin.HandlerFunc {
	return sessionHandler
}
