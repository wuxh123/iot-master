package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/model"
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

	var user model.User
	has, err := db.Engine.Where("username=?", obj.Username).Get(&user)
	if err != nil {
		replyError(ctx, err)
		return
	}

	if !has {
		replyFail(ctx, "无此用户")
		return
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
