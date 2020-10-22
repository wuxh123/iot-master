package api

import (
	"encoding/json"
	"git.zgwit.com/zgwit/iot-admin/conf"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/gorilla/mux"
	"net/http"
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
//
//var (
//	cookieNameForSessionID = "iot-admin"
//	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
//)
//
//func mustLogin(ctx iris.Context) {
//	session := sess.Start(ctx)
//	if user := session.Get("user"); user != nil {
//		ctx.Next()
//	} else {
//		//TODO 检查OAuth2返回的code，进一步获取用户信息，放置到session中
//
//		ctx.StatusCode(iris.StatusUnauthorized)
//		ctx.JSON(iris.Map{"ok": false, "error": "Unauthorized"})
//		//ctx.Abort()
//	}
//}

func RegisterRoutes(app *mux.Router) {

	if conf.Config.SysAdmin.Enable {
		//检查 session，必须登录
		//app.Use(mustLogin)
	} else if conf.Config.BaseAuth.Enable {
		//检查HTTP认证
		//app.Use(gin.BasicAuth(gin.Accounts(conf.Config.BaseAuth.Users)))
		//authConfig := basicauth.Config{
		//	Users: conf.Config.BaseAuth.Users,
		//}
		//app.Use(basicauth.New(authConfig))
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
	app.HandleFunc("/project/:id/tunnels", curdApiListById(mod, "project_id")).Methods("POST")
	app.HandleFunc("/tunnels", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/tunnel", curdApiCreate(mod, nil)).Methods("POST")            //TODO 启动
	app.HandleFunc("/tunnel/:id", curdApiDelete(mod, nil)).Methods("DELETE")      //TODO 停止
	app.HandleFunc("/tunnel/:id", curdApiModify(mod, fields, nil)).Methods("PUT") //TODO 重新启动
	app.HandleFunc("/tunnel/:id", curdApiGet(mod)).Methods("GET")

	app.HandleFunc("/tunnel/:id/start", tunnelStart).Methods("GET")
	app.HandleFunc("/tunnel/:id/stop", tunnelStop).Methods("GET")

	//app.HandleFunc("/channel/:id/links")

	//连接管理
	mod = reflect.TypeOf(models.Link{})
	fields = []string{"name"}
	app.HandleFunc("/tunnel/:id/links", curdApiListById(mod, "tunnel_id")).Methods("POST")
	app.HandleFunc("/links", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/link/:id", curdApiDelete(mod, nil)).Methods("DELETE")    //TODO 停止
	app.HandleFunc("/link/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/link/:id", curdApiGet(mod)).Methods("GET")

	mod = reflect.TypeOf(models.Device{})
	fields = []string{"name"}
	app.HandleFunc("/project/:id/devices", curdApiListById(mod, "project_id")).Methods("POST")
	app.HandleFunc("/devices", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/device", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/device/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/device/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/device/:id", curdApiGet(mod)).Methods("GET")

	mod = reflect.TypeOf(models.Location{})
	fields = []string{"name"}
	app.HandleFunc("/device/:id/locations", curdApiListById(mod, "device_id")).Methods("POST")
	//app.HandleFunc("/locations", curdApiList(mod)).Methods("POST")
	//app.HandleFunc("/location", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/location/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	//app.HandleFunc("/location/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/location/:id", curdApiGet(mod)).Methods("GET")

	//插件管理
	mod = reflect.TypeOf(models.Plugin{})
	fields = []string{"name"}
	app.HandleFunc("/plugins", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/plugin", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/plugin/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/plugin/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/plugin/:id", curdApiGet(mod)).Methods("GET")

	//模型管理
	mod = reflect.TypeOf(models.Project{})
	fields = []string{"name"}
	app.HandleFunc("/projects", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/project", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/project/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/project/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/project/:id", curdApiGet(mod)).Methods("GET")

	//app.HandleFunc("/project/:id/tunnels", nop)
	//app.HandleFunc("/project/:id/variables", nop)
	//app.HandleFunc("/project/:id/batches", nop)
	//app.HandleFunc("/project/:id/jobs", nop)
	//app.HandleFunc("/project/:id/strategies", nop)

	//app.HandleFunc("/project/import", projectImport).Methods("POST")
	//app.HandleFunc("/project/:id/export", projectExport).Methods("GET")
	//app.HandleFunc("/project/:id/deploy", projectDeploy).Methods("GET")

	mod = reflect.TypeOf(models.ProjectElement{})
	fields = []string{"name"}
	app.HandleFunc("/project/:id/elements", curdApiListById(mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/elements", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/project/element", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/project/element/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/project/element/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/project/element/:id", curdApiGet(mod)).Methods("GET")

	mod = reflect.TypeOf(models.ProjectJob{})
	fields = []string{"name"}
	app.HandleFunc("/project/:id/jobs", curdApiListById(mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/jobs", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/project/job", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/project/job/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/project/job/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/project/job/:id", curdApiGet(mod)).Methods("GET")

	mod = reflect.TypeOf(models.ProjectStrategy{})
	fields = []string{"name"}
	app.HandleFunc("/project/:id/strategies", curdApiListById(mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/strategies", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/project/strategy", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/project/strategy/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/project/strategy/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/project/strategy/:id", curdApiGet(mod)).Methods("GET")

	//元件管理
	mod = reflect.TypeOf(models.Element{})
	fields = []string{"name"}
	app.HandleFunc("/elements", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/element", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/element/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/element/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/element/:id", curdApiGet(mod)).Methods("GET")

	//元件变量
	mod = reflect.TypeOf(models.ElementVariable{})
	fields = []string{"name"}
	app.HandleFunc("/element/:id/variables", curdApiListById(mod, "element_id")).Methods("POST")
	//app.HandleFunc("/element/variables", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/element/variable", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/element/variable/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/element/variable/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/element/variable/:id", curdApiGet(mod)).Methods("GET")

	//元件批量操作
	mod = reflect.TypeOf(models.ElementBatch{})
	fields = []string{"name"}
	app.HandleFunc("/element/:id/batches", curdApiListById(mod, "element_id")).Methods("POST")
	//app.HandleFunc("/element/batches", curdApiList(mod)).Methods("POST")
	app.HandleFunc("/element/batch", curdApiCreate(mod, nil)).Methods("POST")
	app.HandleFunc("/element/batch/:id", curdApiDelete(mod, nil)).Methods("DELETE")
	app.HandleFunc("/element/batch/:id", curdApiModify(mod, fields, nil)).Methods("PUT")
	app.HandleFunc("/element/batch/:id", curdApiGet(mod)).Methods("GET")

}

type Reply struct {
	Ok    bool        `json:"ok"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Total int64       `json:"total,omitempty"`
}

func replyList(writer http.ResponseWriter, data interface{}, total int64) {
	r := Reply{
		Ok:    true,
		Data:  data,
		Total: total,
	}
	b, _ := json.Marshal(r)
	_, _ = writer.Write(b)
}

func replyOk(writer http.ResponseWriter, data interface{}) {
	r := Reply{
		Ok:   true,
		Data: data,
	}
	b, _ := json.Marshal(r)
	_, _ = writer.Write(b)
}

func replyFail(writer http.ResponseWriter, err string) {
	r := Reply{
		Error: err,
	}
	b, _ := json.Marshal(r)
	_, _ = writer.Write(b)
}

func replyError(writer http.ResponseWriter, err error) {
	r := Reply{
		Error: err.Error(),
	}
	b, _ := json.Marshal(r)
	_, _ = writer.Write(b)
}

func nop(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusForbidden)
	_, _ = writer.Write([]byte("Unsupported"))
}
