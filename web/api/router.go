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
	mod = reflect.TypeOf(models.Project{})
	app.POST("/projects", curdApiList(mod))
	app.POST("/project", curdApiCreate(mod, nil))
	app.DELETE("/project/:id", curdApiDelete(mod, nil))
	app.PUT("/project/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/project/:id", curdApiGet(mod))

	//app.GET("/project/:id/tunnels", nop)
	//app.GET("/project/:id/variables", nop)
	//app.GET("/project/:id/batches", nop)
	//app.GET("/project/:id/jobs", nop)
	//app.GET("/project/:id/strategies", nop)

	app.POST("/project/import", projectImport)
	app.GET("/project/:id/export", projectExport)

	app.GET("/project/:id/deploy", projectDeploy)

	mod = reflect.TypeOf(models.ProjectAdapter{})
	app.POST("/project-adapters", curdApiList(mod))
	app.POST("/project-adapter", curdApiCreate(mod, nil))
	app.DELETE("/project-adapter/:id", curdApiDelete(mod, nil))
	app.PUT("/project-adapter/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/project-adapter/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ProjectVariable{})
	app.POST("/project-variables", curdApiList(mod))
	app.POST("/project-variable", curdApiCreate(mod, nil))
	app.DELETE("/project-variable/:id", curdApiDelete(mod, nil))
	app.PUT("/project-variable/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/project-variable/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ProjectBatch{})
	app.POST("/project-batches", curdApiList(mod))
	app.POST("/project-batch", curdApiCreate(mod, nil))
	app.DELETE("/project-batch/:id", curdApiDelete(mod, nil))
	app.PUT("/project-batch/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/project-batch/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ProjectJob{})
	app.POST("/project-jobs", curdApiList(mod))
	app.POST("/project-job", curdApiCreate(mod, nil))
	app.DELETE("/project-job/:id", curdApiDelete(mod, nil))
	app.PUT("/project-job/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/project-job/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ProjectStrategy{})
	app.POST("/project-strategies", curdApiList(mod))
	app.POST("/project-strategy", curdApiCreate(mod, nil))
	app.DELETE("/project-strategy/:id", curdApiDelete(mod, nil))
	app.PUT("/project-strategy/:id", curdApiModify(mod, []string{}, nil))
	app.GET("/project-strategy/:id", curdApiGet(mod))
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
