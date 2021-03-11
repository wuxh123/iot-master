package api

import (
	"mydtu/db"
	"mydtu/model"
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

	 err := db.DB("user").One("ID", u.Id, u)
	if err != nil {
		//if err == storm.ErrNotFound {
		//	replyError(c, err)
		//	return
		//}
		replyError(c, err)
		return
	}


	if u.Password != obj.Origin {
		replyFail(c, "密钥错误")
		return
	}

	//u.Password = obj.Password
	err = db.DB("user").UpdateField(u, "Password", obj.Password)
	if err != nil {
		replyError(c, err)
		return
	}

	session.Clear()
	_ = session.Save()

	replyOk(c, nil)
}
