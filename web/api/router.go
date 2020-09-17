package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"git.zgwit.com/iot/dtu-admin/conf"
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
	Id int `uri:"id"`
}

type paramId2 struct {
	Id  int `uri:"id"`
	Id2 int `uri:"id2"`
}

func mustLogin(c *gin.Context) {
	session := sessions.Default(c)
	if user := session.Get("user"); user != nil {
		c.Next()
	} else {
		//TODO 检查OAuth2返回的code，进一步获取用户信息，放置到session中

		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "Unauthorized"})
		c.Abort()
	}
}

func RegisterRoutes(app *gin.RouterGroup) {

	if conf.Config.SysAdmin.Enable {
		//检查 session，必须登录
		app.Use(mustLogin)

		app.GET("/mqtt", mqtt)
	} else if conf.Config.BaseAuth.Enable {
		app.GET("/mqtt", mqtt)

		//检查HTTP认证
		app.Use(gin.BasicAuth(gin.Accounts(conf.Config.BaseAuth.Users)))
	} else {
		//支持匿名访问
	}

	//TODO 转移至子目录，并使用中间件，检查session及权限
	app.POST("/channels", channels)
	app.POST("/channel", channelCreate)
	app.DELETE("/channel/:id", channelDelete)
	app.PUT("/channel/:id", channelModify)
	app.GET("/channel/:id", channelGet)
	app.GET("/channel/:id/start", channelStart)
	app.GET("/channel/:id/stop", channelStop)

	app.POST("/channel/:id/links")

	app.POST("/links", links)
	app.DELETE("/link/:id", linkDelete)
	app.PUT("/link/:id", linkModify)
	app.GET("/link/:id", linkGet)

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
