package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/dtu-admin/storage"
	"github.com/zgwit/dtu-admin/types"
	"time"
)

type loginObj struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func authLogin(ctx *gin.Context) {
	session := sessions.Default(ctx)

	var obj loginObj
	err := ctx.ShouldBindJSON(&obj)
	if err != nil {
		replyFail(ctx, "参数解析错误"+err.Error())
		return
	}

	userDB := storage.DB("user")
	var user types.User
	err = userDB.One("username", obj.Username, &user)
	if err != nil {
		//初始化root 账户
		if obj.Username == "admin" {
			user.Username = "admin"
			user.Password = "123456"
			user.Created = time.Now()
			userDB.Save(&user)
		} else {
			replyFail(ctx, "无此用户")
			return
		}
	}

	if obj.Password != user.Password {
		replyFail(ctx, "密码错误")
		return
	}

	session.Set("user", user)
	_ = session.Save()

	replyOk(ctx, user)
}

func authLogout(ctx *gin.Context) {

}

func authPassword(ctx *gin.Context) {

}
