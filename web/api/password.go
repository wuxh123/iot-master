package api

import (
	"iot-master/db"
	"iot-master/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type passwordObj struct {
	Origin   string `json:"origin"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func password(c *gin.Context) {
	session := sessions.Default(c)

	var obj passwordObj
	if err := c.ShouldBindJSON(&obj); err != nil {
		replyError(c, err)
		return
	}

	user := session.Get("user")
	//log.Println("user", user)
	u := user.(*model.User)
	//log.Println("u", u)

	has, err := db.Engine.ID(u.Id).Get(u)
	if !has {
		replyError(c, err)
		return
	} else if err != nil {
		replyError(c, err)
		return
	}

	if u.Password != obj.Origin {
		replyFail(c, "密钥错误")
		return
	}

	u.Password = obj.Password
	_, err = db.Engine.ID(u.Id).Cols("password").Update(u)
	if err != nil {
		replyError(c, err)
		return
	}

	session.Clear()
	_ = session.Save()

	replyOk(c, nil)
}
