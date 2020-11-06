package api

import (
	"encoding/json"
	"git.zgwit.com/zgwit/iot-admin/conf"
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/gorilla/mux"
	"github.com/quasoft/memstore"
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

func RegisterRoutes(app *mux.Router) {

	if conf.Config.SysAdmin.Enable {
		//启用session
		store := memstore.NewMemStore([]byte("iot-admin"), []byte("iot-admin"))
		app.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				sess, err := store.Get(request, "iot-admin")
				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}
				if sess.IsNew {
					_ = sess.Save(request, writer)
				}
				//TODO 检查session，及权限
				next.ServeHTTP(writer, request)
			})
		})
		//检查 session，必须登录
		//app.Use(mustLogin)
	} else if conf.Config.BaseAuth.Enable {
		//检查HTTP认证
		app.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if username, password, ok := request.BasicAuth(); ok {
					if pass, ok := conf.Config.BaseAuth.Users[username]; ok && password == pass {
						next.ServeHTTP(writer, request)
						return
					}
				}

				writer.Header().Set("WWW-Authenticate", `Basic realm="请输入用户名和密码"`)
				http.Error(writer, "Unauthorised", http.StatusUnauthorized)
			})
		})
		//app.Use(gin.BasicAuth(gin.Accounts(conf.Config.BaseAuth.Users)))
	} else {
		//支持匿名访问
	}

	//TODO 转移至子目录，并使用中间件，检查session及权限
	mod := reflect.TypeOf(models.Tunnel{})
	store := db.DB("tunnel")
	app.HandleFunc("/project/{id}/tunnels", curdApiListById(store, mod, "project_id")).Methods("POST")
	app.HandleFunc("/tunnels", curdApiList(store, mod)).Methods("POST")
	app.HandleFunc("/tunnel", curdApiCreate(store, mod, nil, nil)).Methods("POST")             //TODO 启动
	app.HandleFunc("/tunnel/{id}", curdApiDelete(store, mod, nil, nil)).Methods("DELETE")      //TODO 停止
	app.HandleFunc("/tunnel/{id}", curdApiModify(store, mod, nil, nil)).Methods("PUT", "POST") //TODO 重新启动
	app.HandleFunc("/tunnel/{id}", curdApiGet(store, mod)).Methods("GET")

	app.HandleFunc("/tunnel/{id}/start", tunnelStart).Methods("GET")
	app.HandleFunc("/tunnel/{id}/stop", tunnelStop).Methods("GET")

	//app.HandleFunc("/channel/{id}/links")

	//连接管理
	mod = reflect.TypeOf(models.Link{})
	store = db.DB("link")
	app.HandleFunc("/tunnel/{id}/links", curdApiListById(store, mod, "tunnel_id")).Methods("POST")
	app.HandleFunc("/links", curdApiList(store, mod)).Methods("POST")
	app.HandleFunc("/link/{id}", curdApiDelete(store, mod, nil, nil)).Methods("DELETE") //TODO 停止
	app.HandleFunc("/link/{id}", curdApiModify(store, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/link/{id}", curdApiGet(store, mod)).Methods("GET")

	//设备管理
	mod = reflect.TypeOf(models.Device{})
	store = db.DB("device")
	app.HandleFunc("/project/{id}/devices", curdApiListById(store, mod, "project_id")).Methods("POST")
	app.HandleFunc("/devices", curdApiList(store, mod)).Methods("POST")
	app.HandleFunc("/device", curdApiCreate(store, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/device/{id}", curdApiDelete(store, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/device/{id}", curdApiModify(store, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/device/{id}", curdApiGet(store, mod)).Methods("GET")

	//插件管理
	mod = reflect.TypeOf(models.Plugin{})
	store = db.DB("plugin")
	app.HandleFunc("/plugins", curdApiList(store, mod)).Methods("POST")
	app.HandleFunc("/plugin", curdApiCreate(store, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/plugin/{id}", curdApiDelete(store, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/plugin/{id}", curdApiModify(store, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/plugin/{id}", curdApiGet(store, mod)).Methods("GET")

	//项目管理
	mod = reflect.TypeOf(models.Project{})
	store = db.DB("project")
	app.HandleFunc("/projects", curdApiList(store, mod)).Methods("POST")
	app.HandleFunc("/project", curdApiCreate(store, mod, projectBeforeCreate, projectAfterCreate)).Methods("POST")
	app.HandleFunc("/project/{id}", curdApiDelete(store, mod, nil, projectAfterDelete)).Methods("DELETE")
	app.HandleFunc("/project/{id}", curdApiModify(store, mod, nil, projectAfterModify)).Methods("PUT", "POST")
	app.HandleFunc("/project/{id}", curdApiGet(store, mod)).Methods("GET")

	//app.HandleFunc("/project/import", projectImport).Methods("POST")
	//app.HandleFunc("/project/{id}/export", projectExport).Methods("GET")
	//app.HandleFunc("/project/{id}/deploy", projectDeploy).Methods("GET")

	//项目链接
	mod = reflect.TypeOf(models.ProjectLink{})
	node := store.From("link")
	app.HandleFunc("/project/{id}/links", curdApiListById(node, mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/links", curdApiList(node,mod)).Methods("POST")
	app.HandleFunc("/project/link", curdApiCreate(node, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/project/link/{id}", curdApiDelete(node, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/project/link/{id}", curdApiModify(node, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/project/link/{id}", curdApiGet(node, mod)).Methods("GET")

	//项目元件
	mod = reflect.TypeOf(models.ProjectElement{})
	node = store.From("element")
	app.HandleFunc("/project/{id}/elements", curdApiListById(node, mod, "project_id")).Methods("POST")
	app.HandleFunc("/project/link/{id}/elements", curdApiListById(node, mod, "link_id")).Methods("POST")
	//app.HandleFunc("/project/elements", curdApiList(node,mod)).Methods("POST")
	app.HandleFunc("/project/element", curdApiCreate(node, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/project/element/{id}", curdApiDelete(node, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/project/element/{id}", curdApiModify(node, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/project/element/{id}", curdApiGet(node, mod)).Methods("GET")

	//项目变量
	mod = reflect.TypeOf(models.ProjectVariable{})
	node = store.From("variable")
	app.HandleFunc("/project/{id}/variables", curdApiListById(node, mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/variables", curdApiList(node,mod)).Methods("POST")
	app.HandleFunc("/project/variable", curdApiCreate(node, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/project/variable/{id}", curdApiDelete(node, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/project/variable/{id}", curdApiModify(node, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/project/variable/{id}", curdApiGet(node, mod)).Methods("GET")

	//项目检查
	mod = reflect.TypeOf(models.ProjectValidator{})
	node = store.From("validator")
	app.HandleFunc("/project/{id}/validators", curdApiListById(node, mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/validators", curdApiList(node,mod)).Methods("POST")
	app.HandleFunc("/project/validator", curdApiCreate(node, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/project/validator/{id}", curdApiDelete(node, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/project/validator/{id}", curdApiModify(node, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/project/validator/{id}", curdApiGet(node, mod)).Methods("GET")

	//项目功能
	mod = reflect.TypeOf(models.ProjectFunction{})
	node = store.From("function")
	app.HandleFunc("/project/{id}/functions", curdApiListById(node, mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/functions", curdApiList(node,mod)).Methods("POST")
	app.HandleFunc("/project/function", curdApiCreate(node, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/project/function/{id}", curdApiDelete(node, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/project/function/{id}", curdApiModify(node, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/project/function/{id}", curdApiGet(node, mod)).Methods("GET")

	//项目任务
	mod = reflect.TypeOf(models.ProjectJob{})
	node = store.From("job")
	app.HandleFunc("/project/{id}/jobs", curdApiListById(node, mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/jobs", curdApiList(node,mod)).Methods("POST")
	app.HandleFunc("/project/job", curdApiCreate(node, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/project/job/{id}", curdApiDelete(node, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/project/job/{id}", curdApiModify(node, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/project/job/{id}", curdApiGet(node, mod)).Methods("GET")

	//项目策略
	mod = reflect.TypeOf(models.ProjectStrategy{})
	node = store.From("strategy")
	app.HandleFunc("/project/{id}/strategies", curdApiListById(node, mod, "project_id")).Methods("POST")
	//app.HandleFunc("/project/strategies", curdApiList(node,mod)).Methods("POST")
	app.HandleFunc("/project/strategy", curdApiCreate(node, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/project/strategy/{id}", curdApiDelete(node, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/project/strategy/{id}", curdApiModify(node, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/project/strategy/{id}", curdApiGet(node, mod)).Methods("GET")

	//元件管理
	mod = reflect.TypeOf(models.Element{})
	store = db.DB("element")
	app.HandleFunc("/elements", curdApiList(store, mod)).Methods("POST")
	app.HandleFunc("/element", curdApiCreate(store, mod, elementBeforeCreate, nil)).Methods("POST")
	app.HandleFunc("/element/{id}", curdApiDelete(store, mod, elementBeforeDelete, nil)).Methods("DELETE")
	app.HandleFunc("/element/{id}", curdApiModify(store, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/element/{id}", curdApiGet(store, mod)).Methods("GET")

	//元件变量
	mod = reflect.TypeOf(models.ElementVariable{})
	node = store.From("variable")
	app.HandleFunc("/element/{id}/variables", curdApiListById(node, mod, "project_id")).Methods("POST")
	//app.HandleFunc("/element/variables", curdApiList(node,mod)).Methods("POST")
	app.HandleFunc("/element/variable", curdApiCreate(node, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/element/variable/{id}", curdApiDelete(node, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/element/variable/{id}", curdApiModify(node, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/element/variable/{id}", curdApiGet(node, mod)).Methods("GET")

}

type Reply struct {
	Ok    bool        `json:"ok"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
type ReplyList struct {
	Ok    bool        `json:"ok"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Total int         `json:"total"`
}

func replyList(writer http.ResponseWriter, data interface{}, total int) {
	r := ReplyList{
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
