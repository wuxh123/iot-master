package web

import (
	"git.zgwit.com/zgwit/iot-admin/conf"
	"git.zgwit.com/zgwit/iot-admin/web/api"
	"git.zgwit.com/zgwit/iot-admin/web/open"
	wwwFiles "git.zgwit.com/zgwit/iot-admin/web/www"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func Serve() {
	if !conf.Config.Web.Debug {
		//gin.SetMode(gin.ReleaseMode)

	}

	app := mux.NewRouter()

	//加入swagger会增加10MB多体积，使用github.com/zgwit/swagger-files，去除Map文件，可以节省7MB左右
	//Swagger文档，需要先执行swag init生成文档
	//app.Get("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//MQTT
	app.HandleFunc("/mqtt", mqtt).Methods("GET")
	//透传
	app.HandleFunc("/peer", peer).Methods("GET")

	//开放接口
	open.RegisterRoutes(app.PathPrefix("/open").Subrouter())

	//授权检查，启用了SysAdmin的OAuth2，就不能再使用基本HTTP认证了
	//if conf.Config.SysAdmin.Enable {
	//	//注册OAuth2相关接口
	//	RegisterOauthRoutes(app)
	//}


	//注册前端接口
	api.RegisterRoutes(app.PathPrefix("/api").Subrouter())

	//未登录，访问前端文件，跳转到OAuth2登录
	if conf.Config.SysAdmin.Enable {
		//app.Use(func(c iris.Context) {
		//	//session := sessions.Default(c)
		//	//if user := session.Get("user"); user != nil {
		//	//	c.Next()
		//	//} else {
		//	//	//TODO 拼接 OAuth2链接，需要AppKey和Secret
		//	//	url := conf.Config.SysAdmin.Addr + "?redirect_uri="
		//	//	c.Redirect(http.StatusFound, url)
		//	//}
		//})
	} else if conf.Config.BaseAuth.Enable {
		//开启基本HTTP认证
		//app.Use(gin.BasicAuth(gin.Accounts(conf.Config.BaseAuth.Users)))
	}

	//前端静态文件
	//app.PathPrefix("/").Handler(wwwFiles.Handler).Methods("GET")
	app.PathPrefix("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//支持前端框架的无“#”路由
		if request.RequestURI == "/" {
			request.URL.Path = "index.html"
		} else if _, err := wwwFiles.FS.Stat(wwwFiles.CTX, request.RequestURI); err != nil {
			request.URL.Path = "index.html"
		}
		//文件失效期已经在Handler中处理
		wwwFiles.Handler.ServeHTTP(writer, request)
	}).Methods("GET")

	//监听HTTP
	srv := &http.Server{
		Addr:         conf.Config.Web.Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: app,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("HTTP 服务启动错误", err)
	}
}
