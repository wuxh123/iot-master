package api

import (
	"git.zgwit.com/zgwit/iot-admin/conf"
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
	mod := reflect.TypeOf(models.Tunnel{})
	app.POST("/tunnels", curdApiList(mod))
	//app.POST("/tunnel", curdApiCreate(mod, nil))
	app.DELETE("/tunnel/:id", curdApiDelete(mod, nil))
	app.PUT("/tunnel/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/tunnel/:id", curdApiGet(mod))

	app.GET("/tunnel/:id/start", tunnelStart)
	app.GET("/tunnel/:id/stop", tunnelStop)

	//app.POST("/channel/:id/links")

	//连接管理
	mod = reflect.TypeOf(models.Link{})
	app.POST("/links", curdApiList(mod))
	app.DELETE("/link/:id", curdApiDelete(mod, nil))
	app.PUT("/link/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/link/:id", curdApiGet(mod))

	//插件管理
	mod = reflect.TypeOf(models.Plugin{})
	app.POST("/plugins", curdApiList(mod))
	app.POST("/plugin", curdApiCreate(mod, nil))
	app.DELETE("/plugin/:id", curdApiDelete(mod, nil))
	app.PUT("/plugin/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/plugin/:id", curdApiGet(mod))

	//模型管理
	mod = reflect.TypeOf(models.Model{})
	app.POST("/models", curdApiList(mod))
	app.POST("/model", curdApiCreate(mod, nil))
	app.DELETE("/model/:id", curdApiDelete(mod, nil))
	app.PUT("/model/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/model/:id", curdApiGet(mod))

	//app.GET("/model/:id/tunnels", nop)
	//app.GET("/model/:id/variables", nop)
	//app.GET("/model/:id/batches", nop)
	//app.GET("/model/:id/jobs", nop)
	//app.GET("/model/:id/strategies", nop)

	app.POST("/model/import", modelImport)
	app.GET("/model/:id/export", modelExport)

	app.GET("/model/:id/refresh", modelRefresh)

	mod = reflect.TypeOf(models.ModelAdapter{})
	app.POST("/model-adapters", curdApiList(mod))
	app.POST("/model-adapter", curdApiCreate(mod, nil))
	app.DELETE("/model-adapter/:id", curdApiDelete(mod, nil))
	app.PUT("/model-adapter/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/model-adapter/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ModelVariable{})
	app.POST("/model-variables", curdApiList(mod))
	app.POST("/model-variable", curdApiCreate(mod, nil))
	app.DELETE("/model-variable/:id", curdApiDelete(mod, nil))
	app.PUT("/model-variable/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/model-variable/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ModelBatch{})
	app.POST("/model-batches", curdApiList(mod))
	app.POST("/model-batch", curdApiCreate(mod, nil))
	app.DELETE("/model-batch/:id", curdApiDelete(mod, nil))
	app.PUT("/model-batch/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/model-batch/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ModelJob{})
	app.POST("/model-jobs", curdApiList(mod))
	app.POST("/model-job", curdApiCreate(mod, nil))
	app.DELETE("/model-job/:id", curdApiDelete(mod, nil))
	app.PUT("/model-job/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/model-job/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ModelStrategy{})
	app.POST("/model-strategies", curdApiList(mod))
	app.POST("/model-strategy", curdApiCreate(mod, nil))
	app.DELETE("/model-strategy/:id", curdApiDelete(mod, nil))
	app.PUT("/model-strategy/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/model-strategy/:id", curdApiGet(mod))
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
