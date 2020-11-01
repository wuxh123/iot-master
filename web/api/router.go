package api

import (
	"encoding/json"
	"git.zgwit.com/zgwit/iot-admin/conf"
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
	model := "tunnel"
	app.HandleFunc("/project/{id}/tunnels", curdApiListById(model, mod, "project_id")).Methods("POST")
	app.HandleFunc("/tunnels", curdApiList(model, mod)).Methods("POST")
	app.HandleFunc("/tunnel", curdApiCreate(model, mod, nil, nil)).Methods("POST")             //TODO 启动
	app.HandleFunc("/tunnel/{id}", curdApiDelete(model, mod, nil, nil)).Methods("DELETE")      //TODO 停止
	app.HandleFunc("/tunnel/{id}", curdApiModify(model, mod, nil, nil)).Methods("PUT", "POST") //TODO 重新启动
	app.HandleFunc("/tunnel/{id}", curdApiGet(model, mod)).Methods("GET")

	app.HandleFunc("/tunnel/{id}/start", tunnelStart).Methods("GET")
	app.HandleFunc("/tunnel/{id}/stop", tunnelStop).Methods("GET")

	//app.HandleFunc("/channel/{id}/links")

	//连接管理
	mod = reflect.TypeOf(models.Link{})
	model = "link"
	app.HandleFunc("/tunnel/{id}/links", curdApiListById(model, mod, "tunnel_id")).Methods("POST")
	app.HandleFunc("/links", curdApiList(model, mod)).Methods("POST")
	app.HandleFunc("/link/{id}", curdApiDelete(model, mod, nil, nil)).Methods("DELETE") //TODO 停止
	app.HandleFunc("/link/{id}", curdApiModify(model, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/link/{id}", curdApiGet(model, mod)).Methods("GET")

	//插件管理
	mod = reflect.TypeOf(models.Plugin{})
	model = "plugin"
	app.HandleFunc("/plugins", curdApiList(model, mod)).Methods("POST")
	app.HandleFunc("/plugin", curdApiCreate(model, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/plugin/{id}", curdApiDelete(model, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/plugin/{id}", curdApiModify(model, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/plugin/{id}", curdApiGet(model, mod)).Methods("GET")

	//项目管理
	mod = reflect.TypeOf(models.Project{})
	model = "project"
	app.HandleFunc("/projects", curdApiList(model, mod)).Methods("POST")
	app.HandleFunc("/project", curdApiCreate(model, mod, projectBeforeCreate, projectAfterCreate)).Methods("POST")
	app.HandleFunc("/project/{id}", curdApiDelete(model, mod, nil, projectAfterDelete)).Methods("DELETE")
	app.HandleFunc("/project/{id}", curdApiModify(model, mod, nil, projectAfterModify)).Methods("PUT", "POST")
	app.HandleFunc("/project/{id}", curdApiGet(model, mod)).Methods("GET")

	//app.HandleFunc("/project/import", projectImport).Methods("POST")
	//app.HandleFunc("/project/{id}/export", projectExport).Methods("GET")
	//app.HandleFunc("/project/{id}/deploy", projectDeploy).Methods("GET")


	//模板管理
	mod = reflect.TypeOf(models.ProjectTemplate{})
	model = "template"
	app.HandleFunc("/templates", curdApiList(model, mod)).Methods("POST")
	app.HandleFunc("/template", curdApiCreate(model, mod, nil, nil)).Methods("POST")
	app.HandleFunc("/template/{id}", curdApiDelete(model, mod, nil, nil)).Methods("DELETE")
	app.HandleFunc("/template/{id}", curdApiModify(model, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/template/{id}", curdApiGet(model, mod)).Methods("GET")

	//元件管理
	mod = reflect.TypeOf(models.Element{})
	model = "element"
	app.HandleFunc("/elements", curdApiList(model, mod)).Methods("POST")
	app.HandleFunc("/element", curdApiCreate(model, mod, elementBeforeCreate, nil)).Methods("POST")
	app.HandleFunc("/element/{id}", curdApiDelete(model, mod, elementBeforeDelete, nil)).Methods("DELETE")
	app.HandleFunc("/element/{id}", curdApiModify(model, mod, nil, nil)).Methods("PUT", "POST")
	app.HandleFunc("/element/{id}", curdApiGet(model, mod)).Methods("GET")
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
