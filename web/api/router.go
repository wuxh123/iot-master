package api

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/dtu-admin/model"
	"net/http"
)

type paramFilter struct {
	Key   string   `form:"key"`
	Value []string `form:"value"`
}

type paramSearch struct {
	Offset    int           `form:"offset"`
	Length    int           `form:"length"`
	SortKey   string        `form:"sortKey"`
	SortOrder string        `form:"sortOrder"`
	Filters   []paramFilter `form:"filters"`
	Keyword   string        `form:"keyword"`
}

type paramId struct {
	Id int64 `uri:"id"`
}

type paramId2 struct {
	Id  int64 `uri:"id"`
	Id2 int64 `uri:"id2"`
}

func RegisterRoutes(app *gin.RouterGroup) {
	//注册 User类型
	gob.Register(&model.User{})
	//启用session
	app.Use(sessions.Sessions("dtu-admin", memstore.NewStore([]byte("dtu-admin-secret"))))

	app.POST("/login", authLogin)

	//检查登录状态
	//app.Use(func(c *gin.Context) {
	//	session := sessions.Default(c)
	//	if user := session.Get("user"); user != nil {
	//		c.Next()
	//	} else {
	//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	//		c.Abort()
	//	}
	//})

	app.DELETE("/logout", authLogout)
	app.POST("/password", authPassword)

	//TODO 转移至子目录，并使用中间件，检查session及权限
	app.POST("/channels", channels)
	app.POST("/channel", channelCreate)
	app.DELETE("/channel/:id", channelDelete)
	app.PUT("/channel/:id", channelModify)
	app.GET("/channel/:id", channelGet)
	app.GET("/channel/:id/start", channelStart)
	app.GET("/channel/:id/stop", channelStop)

	app.GET("/channel/:id/links")
	app.POST("/channel/:id/links")
	app.DELETE("/channel/:id/link/:id2") //关闭连接
	app.GET("/channel/:id/link/:id2/statistic")
	app.GET("/channel/:id/link/:id2/pipe") //转Websocket透传

	app.POST("/links", links)
	app.DELETE("/link/:id", linkDelete)
	app.PUT("/link/:id", linkModify)
	app.GET("/link/:id", linkGet)
	app.GET("/link/:id/monitor", linkMonitor)
	app.POST("/link/:id/send", linkSend)
	//app.GET("/link/:id/watch", linkWatch)
	//app.GET("/link/:id/stop", linkStop)

}

func replyOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": data,
	})
}

func replyFail(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusOK, gin.H{
		"ok":    false,
		"error": err,
	})
}

func replyError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"ok":    false,
		"error": err.Error(),
	})
}


func nop(c *gin.Context) {
	c.String(http.StatusForbidden, "Unsupported")
}

func mustLogin(c *gin.Context) {
	//测试
	session := sessions.Default(c)
	if user := session.Get("user"); user != nil {
		c.Next()
	} else {
		//c.Redirect(http.StatusSeeOther, "/login")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}
