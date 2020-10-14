package api

import (
	"git.zgwit.com/zgwit/iot-admin/internal/conf"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
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
	} else if conf.Config.BaseAuth.Enable {
		//检查HTTP认证
		app.Use(gin.BasicAuth(gin.Accounts(conf.Config.BaseAuth.Users)))
	} else {
		//支持匿名访问
	}

	//TODO 转移至子目录，并使用中间件，检查session及权限
	typ := reflect.TypeOf(models.Tunnel{})
	app.POST("/tunnels", curdApiList(typ))
	//app.POST("/tunnel", curdApiCreate(typ, nil))
	app.DELETE("/tunnel/:id", curdApiDelete(typ, nil))
	app.PUT("/tunnel/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/tunnel/:id", curdApiGet(typ))

	app.GET("/tunnel/:id/start", tunnelStart)
	app.GET("/tunnel/:id/stop", tunnelStop)

	//app.POST("/channel/:id/links")

	//连接管理
	typ = reflect.TypeOf(models.Link{})
	app.POST("/links", curdApiList(typ))
	app.DELETE("/link/:id", curdApiDelete(typ, nil))
	app.PUT("/link/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/link/:id", curdApiGet(typ))

	//插件管理
	typ = reflect.TypeOf(models.Plugin{})
	app.POST("/plugins", curdApiList(typ))
	app.POST("/plugin", curdApiCreate(typ, nil))
	app.DELETE("/plugin/:id", curdApiDelete(typ, nil))
	app.PUT("/plugin/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/plugin/:id", curdApiGet(typ))

	//模型管理
	typ = reflect.TypeOf(models.Model{})
	app.POST("/models", curdApiList(typ))
	app.POST("/model", curdApiCreate(typ, nil))
	app.DELETE("/model/:id", curdApiDelete(typ, nil))
	app.PUT("/model/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/model/:id", curdApiGet(typ))

	//app.GET("/model/:id/tunnels", nop)
	//app.GET("/model/:id/variables", nop)
	//app.GET("/model/:id/batches", nop)
	//app.GET("/model/:id/jobs", nop)
	//app.GET("/model/:id/strategies", nop)

	app.POST("/model/import", modelImport)
	app.GET("/model/:id/export", modelExport)

	app.GET("/model/:id/refresh", modelRefresh)

	typ = reflect.TypeOf(models.ModelTunnel{})
	app.POST("/model-tunnels", curdApiList(typ))
	app.POST("/model-tunnel", curdApiCreate(typ, nil))
	app.DELETE("/model-tunnel/:id", curdApiDelete(typ, nil))
	app.PUT("/model-tunnel/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/model-tunnel/:id", curdApiGet(typ))

	typ = reflect.TypeOf(models.ModelVariable{})
	app.POST("/model-variables", curdApiList(typ))
	app.POST("/model-variable", curdApiCreate(typ, nil))
	app.DELETE("/model-variable/:id", curdApiDelete(typ, nil))
	app.PUT("/model-variable/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/model-variable/:id", curdApiGet(typ))

	typ = reflect.TypeOf(models.ModelBatch{})
	app.POST("/model-batches", curdApiList(typ))
	app.POST("/model-batch", curdApiCreate(typ, nil))
	app.DELETE("/model-batch/:id", curdApiDelete(typ, nil))
	app.PUT("/model-batch/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/model-batch/:id", curdApiGet(typ))

	typ = reflect.TypeOf(models.ModelJob{})
	app.POST("/model-jobs", curdApiList(typ))
	app.POST("/model-job", curdApiCreate(typ, nil))
	app.DELETE("/model-job/:id", curdApiDelete(typ, nil))
	app.PUT("/model-job/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/model-job/:id", curdApiGet(typ))

	typ = reflect.TypeOf(models.ModelStrategy{})
	app.POST("/model-strategies", curdApiList(typ))
	app.POST("/model-strategy", curdApiCreate(typ, nil))
	app.DELETE("/model-strategy/:id", curdApiDelete(typ, nil))
	app.PUT("/model-strategy/:id", curdApiModify(typ, []string{}, nil))
	app.GET("/model-strategy/:id", curdApiGet(typ))
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
