package api

import (
	"git.zgwit.com/zgwit/iot-admin/conf"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/basicauth"
	"github.com/kataras/iris/v12/sessions"
	"reflect"
)

type paramFilter struct {
	Key    string        `form:"key"`
	Values []interface{} `form:"value"`
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

var (
	cookieNameForSessionID ="iot-admin"
	sess = sessions.New(sessions.Config{Cookie:cookieNameForSessionID})
)

func mustLogin(ctx iris.Context) {
	session := sess.Start(ctx)
	if user := session.Get("user"); user != nil {
		ctx.Next()
	} else {
		//TODO 检查OAuth2返回的code，进一步获取用户信息，放置到session中

		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"ok": false, "error": "Unauthorized"})
		//ctx.Abort()
	}
}

func RegisterRoutes(app iris.Party) {

	if conf.Config.SysAdmin.Enable {
		//检查 session，必须登录
		app.Use(mustLogin)
	} else if conf.Config.BaseAuth.Enable {
		//检查HTTP认证
		//app.Use(gin.BasicAuth(gin.Accounts(conf.Config.BaseAuth.Users)))
		authConfig := basicauth.Config{
			Users: conf.Config.BaseAuth.Users,
		}
		app.Use(basicauth.New(authConfig))
	} else {
		//支持匿名访问
	}

	//TODO 转移至子目录，并使用中间件，检查session及权限
	mod := reflect.TypeOf(models.Tunnel{})
	fields := []string{
		"name", "description", "type", "addr", "timeout",
		"register_enable", "register_regex", "register_min", "register_max",
		"heart_beat_enable", "heart_beat_interval", "heart_beat_content", "heart_beat_is_hex",
		"disabled"}
	app.Post("/project/:id/tunnels", curdApiListById(mod, "project_id"))
	app.Post("/tunnels", curdApiList(mod))
	app.Post("/tunnel", curdApiCreate(mod, nil))            //TODO 启动
	app.Delete("/tunnel/:id", curdApiDelete(mod, nil))      //TODO 停止
	app.Put("/tunnel/:id", curdApiModify(mod, fields, nil)) //TODO 重新启动
	app.Get("/tunnel/:id", curdApiGet(mod))

	app.Get("/tunnel/:id/start", tunnelStart)
	app.Get("/tunnel/:id/stop", tunnelStop)

	//app.Post("/channel/:id/links")

	//连接管理
	mod = reflect.TypeOf(models.Link{})
	fields = []string{"name"}
	app.Post("/tunnel/:id/links", curdApiListById(mod, "tunnel_id"))
	app.Post("/links", curdApiList(mod))
	app.Delete("/link/:id", curdApiDelete(mod, nil)) //TODO 停止
	app.Put("/link/:id", curdApiModify(mod, fields, nil))
	app.Get("/link/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.Device{})
	fields = []string{"name"}
	app.Post("/project/:id/devices", curdApiListById(mod, "project_id"))
	app.Post("/devices", curdApiList(mod))
	app.Post("/device", curdApiCreate(mod, nil))
	app.Delete("/device/:id", curdApiDelete(mod, nil))
	app.Put("/device/:id", curdApiModify(mod, fields, nil))
	app.Get("/device/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.Location{})
	fields = []string{"name"}
	app.Post("/device/:id/locations", curdApiListById(mod, "device_id"))
	//app.Post("/locations", curdApiList(mod))
	//app.Post("/location", curdApiCreate(mod, nil))
	app.Delete("/location/:id", curdApiDelete(mod, nil))
	//app.Put("/location/:id", curdApiModify(mod, fields, nil))
	app.Get("/location/:id", curdApiGet(mod))

	//插件管理
	mod = reflect.TypeOf(models.Plugin{})
	fields = []string{"name"}
	app.Post("/plugins", curdApiList(mod))
	app.Post("/plugin", curdApiCreate(mod, nil))
	app.Delete("/plugin/:id", curdApiDelete(mod, nil))
	app.Put("/plugin/:id", curdApiModify(mod, fields, nil))
	app.Get("/plugin/:id", curdApiGet(mod))

	//模型管理
	mod = reflect.TypeOf(models.Project{})
	fields = []string{"name"}
	app.Post("/projects", curdApiList(mod))
	app.Post("/project", curdApiCreate(mod, nil))
	app.Delete("/project/:id", curdApiDelete(mod, nil))
	app.Put("/project/:id", curdApiModify(mod, fields, nil))
	app.Get("/project/:id", curdApiGet(mod))

	//app.Get("/project/:id/tunnels", nop)
	//app.Get("/project/:id/variables", nop)
	//app.Get("/project/:id/batches", nop)
	//app.Get("/project/:id/jobs", nop)
	//app.Get("/project/:id/strategies", nop)

	app.Post("/project/import", projectImport)
	app.Get("/project/:id/export", projectExport)
	app.Get("/project/:id/deploy", projectDeploy)

	
	mod = reflect.TypeOf(models.ProjectElement{})
	fields = []string{"name"}
	app.Post("/project/:id/elements", curdApiListById(mod, "project_id"))
	//app.Post("/project/elements", curdApiList(mod))
	app.Post("/project/element", curdApiCreate(mod, nil))
	app.Delete("/project/element/:id", curdApiDelete(mod, nil))
	app.Put("/project/element/:id", curdApiModify(mod, fields, nil))
	app.Get("/project/element/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ProjectJob{})
	fields = []string{"name"}
	app.Post("/project/:id/jobs", curdApiListById(mod, "project_id"))
	//app.Post("/project/jobs", curdApiList(mod))
	app.Post("/project/job", curdApiCreate(mod, nil))
	app.Delete("/project/job/:id", curdApiDelete(mod, nil))
	app.Put("/project/job/:id", curdApiModify(mod, fields, nil))
	app.Get("/project/job/:id", curdApiGet(mod))

	mod = reflect.TypeOf(models.ProjectStrategy{})
	fields = []string{"name"}
	app.Post("/project/:id/strategies", curdApiListById(mod, "project_id"))
	//app.Post("/project/strategies", curdApiList(mod))
	app.Post("/project/strategy", curdApiCreate(mod, nil))
	app.Delete("/project/strategy/:id", curdApiDelete(mod, nil))
	app.Put("/project/strategy/:id", curdApiModify(mod, fields, nil))
	app.Get("/project/strategy/:id", curdApiGet(mod))

	//元件管理
	mod = reflect.TypeOf(models.Element{})
	fields = []string{"name"}
	app.Post("/elements", curdApiList(mod))
	app.Post("/element", curdApiCreate(mod, nil))
	app.Delete("/element/:id", curdApiDelete(mod, nil))
	app.Put("/element/:id", curdApiModify(mod, fields, nil))
	app.Get("/element/:id", curdApiGet(mod))

	//元件变量
	mod = reflect.TypeOf(models.ElementVariable{})
	fields = []string{"name"}
	app.Post("/element/:id/variables", curdApiListById(mod, "element_id"))
	//app.Post("/element/variables", curdApiList(mod))
	app.Post("/element/variable", curdApiCreate(mod, nil))
	app.Delete("/element/variable/:id", curdApiDelete(mod, nil))
	app.Put("/element/variable/:id", curdApiModify(mod, fields, nil))
	app.Get("/element/variable/:id", curdApiGet(mod))

	//元件批量操作
	mod = reflect.TypeOf(models.ElementBatch{})
	fields = []string{"name"}
	app.Post("/element/:id/batches", curdApiListById(mod, "element_id"))
	//app.Post("/element/batches", curdApiList(mod))
	app.Post("/element/batch", curdApiCreate(mod, nil))
	app.Delete("/element/batch/:id", curdApiDelete(mod, nil))
	app.Put("/element/batch/:id", curdApiModify(mod, fields, nil))
	app.Get("/element/batch/:id", curdApiGet(mod))

}

func replyOk(ctx iris.Context, data interface{}) {
	ctx.JSON(iris.Map{
		"ok":   true,
		"data": data,
	})
}

func replyFail(ctx iris.Context, err string) {
	ctx.JSON(iris.Map{
		"ok":    false,
		"error": err,
	})
}

func replyError(ctx iris.Context, err error) {
	ctx.JSON(iris.Map{
		"ok":    false,
		"error": err.Error(),
	})
}

func nop(ctx iris.Context) {
	ctx.StatusCode(iris.StatusForbidden)
	ctx.WriteString("Unsupported")
}
