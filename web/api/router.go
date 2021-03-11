package api

import (
	"mydtu/db"
	"mydtu/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type paramFilter struct {
	Key    string        `form:"key"`
	Values []interface{} `form:"value"`
}

type paramKeyword struct {
	Key   string `form:"key"`
	Value string `json:"value"`
}

type paramSearch struct {
	Offset    int            `form:"offset"`
	Length    int            `form:"length"`
	SortKey   string         `form:"sortKey"`
	SortOrder string         `form:"sortOrder"`
	Filters   []paramFilter  `form:"filters"`
	Keywords  []paramKeyword `json:"keywords"`
	//Keyword   string        `form:"keyword"`
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
		c.Set("user", user)
		c.Next()
	} else {
		//TODO 检查OAuth2返回的code，进一步获取用户信息，放置到session中

		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "Unauthorized"})
		c.Abort()
	}
}

func RegisterRoutes(app *gin.RouterGroup) {
	app.POST("/login", login)
	app.Any("/logout", logout)

	//检查 session，必须登录
	app.Use(mustLogin)

	//TODO 转移至子目录，并使用中间件，检查session及权限
	mod := reflect.TypeOf(model.Tunnel{})
	store := db.DB("tunnel")
	app.POST("/project/:id/tunnels", curdApiListById(store, mod, "project_id"))
	app.POST("/tunnels", curdApiList(store, mod))
	app.POST("/tunnel", curdApiCreate(store, mod, nil, nil))       //TODO 启动
	app.DELETE("/tunnel/:id", curdApiDelete(store, mod, nil, nil)) //TODO 停止
	app.PUT("/tunnel/:id", curdApiModify(store, mod, nil, nil))    //TODO 重新启动
	app.GET("/tunnel/:id", curdApiGet(store, mod))

	app.GET("/tunnel/:id/start", tunnelStart)
	app.GET("/tunnel/:id/stop", tunnelStop)

	//app.POST("/channel/:id/links")

	//连接管理
	mod = reflect.TypeOf(model.Link{})
	store = db.DB("link")
	app.POST("/tunnel/:id/links", curdApiListById(store, mod, "tunnel_id"))
	app.POST("/links", curdApiList(store, mod))
	app.DELETE("/link/:id", curdApiDelete(store, mod, nil, nil)) //TODO 停止
	app.PUT("/link/:id", curdApiModify(store, mod, nil, nil))
	app.GET("/link/:id", curdApiGet(store, mod))

	//设备管理
	mod = reflect.TypeOf(model.Device{})
	store = db.DB("device")
	app.POST("/project/:id/devices", curdApiListById(store, mod, "project_id"))
	app.POST("/devices", curdApiList(store, mod))
	app.POST("/device", curdApiCreate(store, mod, nil, nil))
	app.DELETE("/device/:id", curdApiDelete(store, mod, nil, nil))
	app.PUT("/device/:id", curdApiModify(store, mod, nil, nil))
	app.GET("/device/:id", curdApiGet(store, mod))

	//插件管理
	mod = reflect.TypeOf(model.Plugin{})
	store = db.DB("plugin")
	app.POST("/plugins", curdApiList(store, mod))
	app.POST("/plugin", curdApiCreate(store, mod, nil, nil))
	app.DELETE("/plugin/:id", curdApiDelete(store, mod, nil, nil))
	app.PUT("/plugin/:id", curdApiModify(store, mod, nil, nil))
	app.GET("/plugin/:id", curdApiGet(store, mod))

	//项目管理
	mod = reflect.TypeOf(model.Project{})
	store = db.DB("project")
	app.POST("/projects", curdApiList(store, mod))
	app.POST("/project", curdApiCreate(store, mod,  projectBeforeCreate, projectAfterCreate))
	app.DELETE("/project/:id", curdApiDelete(store, mod, nil,projectAfterDelete))
	app.PUT("/project/:id", curdApiModify(store, mod, nil, projectAfterModify))
	app.GET("/project/:id", curdApiGet(store, mod))

	app.POST("/project-import", projectImport)
	app.GET("/project/:id/export", projectExport)
	app.GET("/project/:id/deploy", projectDeploy)

	//元件管理
	mod = reflect.TypeOf(model.Element{})
	store = db.DB("element")
	app.POST("/elements", curdApiList(store, mod))
	app.POST("/element", curdApiCreate(store, mod, elementBeforeCreate,nil))
	app.DELETE("/element/:id", curdApiDelete(store, mod, elementBeforeDelete, nil))
	app.PUT("/element/:id", curdApiModify(store, mod, nil,nil))
	app.GET("/element/:id", curdApiGet(store, mod))
}


func replyList(ctx *gin.Context, data interface{}, total int) {
	ctx.JSON(http.StatusOK, gin.H{"data": data, "total": total})
}

func replyOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func replyFail(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusOK, gin.H{"error": err})
}

func replyError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
}

func nop(ctx *gin.Context) {
	ctx.String(http.StatusForbidden, "Unsupported")
}
