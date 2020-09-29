package api

import (
	"git.zgwit.com/zgwit/iot-admin/internal/conf"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

	//app.POST("/channel/:id/links")

	//连接管理
	app.POST("/links", links)
	app.DELETE("/link/:id", linkDelete)
	app.PUT("/link/:id", linkModify)
	app.GET("/link/:id", linkGet)

	//插件管理
	app.POST("/plugins", plugins)
	app.POST("/plugin", pluginCreate)
	app.DELETE("/plugin/:id", pluginDelete)
	app.PUT("/plugin/:id", pluginModify)
	app.GET("/plugin/:id", pluginGet)

	//模型管理
	app.POST("/models", nop)
	app.POST("/model", nop)
	app.DELETE("/model/:id", nop)
	app.PUT("/model/:id", nop)
	app.GET("/model/:id", nop)

	app.POST("/model/:id/tunnels", nop)
	app.POST("/model/:id/variables", nop)
	app.POST("/model/:id/batches", nop)
	app.POST("/model/:id/jobs", nop)
	app.POST("/model/:id/strategies", nop)

	app.POST("/model-import", modelImport)
	app.GET("/model-export/:id", modelExport)

	app.POST("/tunnels", nop)
	app.POST("/tunnel", nop)
	app.DELETE("/tunnel/:id", nop)
	app.PUT("/tunnel/:id", nop)
	app.GET("/tunnel/:id", nop)

	app.POST("/variables", nop)
	app.POST("/variable", nop)
	app.DELETE("/variable/:id", nop)
	app.PUT("/variable/:id", nop)
	app.GET("/variable/:id", nop)

	app.POST("/batches", nop)
	app.POST("/batch", nop)
	app.DELETE("/batch/:id", nop)
	app.PUT("/batch/:id", nop)
	app.GET("/batch/:id", nop)

	app.POST("/jobs", nop)
	app.POST("/job", nop)
	app.DELETE("/job/:id", nop)
	app.PUT("/job/:id", nop)
	app.GET("/job/:id", nop)

	app.POST("/strategies", nop)
	app.POST("/strategy", nop)
	app.DELETE("/strategy/:id", nop)
	app.PUT("/strategy/:id", nop)
	app.GET("/strategy/:id", nop)

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
