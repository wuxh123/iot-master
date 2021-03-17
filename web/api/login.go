package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"iot-master/db"
	"iot-master/model"
	"net/http"
)

type loginObj struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	var obj loginObj
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		replyError(ctx, err)
		return
	}

	var user model.User
	has, err := db.Engine.Where("username=?", obj.Username).Get(&user)

	if !has {
		replyFail(ctx, "找不到用户")
		return
	} else	if err != nil {
		replyError(ctx, err)
		return
	}
    if user.Password != obj.Password {
        replyFail(ctx, "密钥错误")
        return
    }

	if user.Disabled {
		replyFail(ctx, "用户已禁用")
		return
	}

	//存入session
	session.Set("user", user)
	_ = session.Save()

	replyOk(ctx, user)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	c.JSON(http.StatusOK, gin.H{})
}
