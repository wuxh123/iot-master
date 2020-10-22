package web

import (
	"git.zgwit.com/zgwit/iot-admin/conf"
	"git.zgwit.com/zgwit/iot-admin/web/api"
	"git.zgwit.com/zgwit/iot-admin/web/open"
	wwwFiles "git.zgwit.com/zgwit/iot-admin/web/www"
	"github.com/kataras/iris/v12"
	"log"
	"net/http"
)

func Serve() {
	if !conf.Config.Web.Debug {
		//gin.SetMode(gin.ReleaseMode)
	}


	//GIN初始化
	//app := gin.Default()
	app := iris.New()

	//加入swagger会增加10MB多体积，使用github.com/zgwit/swagger-files，去除Map文件，可以节省7MB左右
	//Swagger文档，需要先执行swag init生成文档
	//app.Get("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//MQTT
	app.Get("/mqtt", mqtt)
	//透传
	app.Get("/peer", peer)

	//开放接口
	open.RegisterRoutes(app.Party("/open"))

	//启用session
	//app.Use(sessions.Sessions("core-admin", memstore.NewStore([]byte("core-admin-secret"))))

	//授权检查，启用了SysAdmin的OAuth2，就不能再使用基本HTTP认证了
	//if conf.Config.SysAdmin.Enable {
	//	//注册OAuth2相关接口
	//	RegisterOauthRoutes(app)
	//}


	//注册前端接口
	api.RegisterRoutes(app.Party("/api"))

	//未登录，访问前端文件，跳转到OAuth2登录
	if conf.Config.SysAdmin.Enable {
		app.Use(func(c iris.Context) {
			//session := sessions.Default(c)
			//if user := session.Get("user"); user != nil {
			//	c.Next()
			//} else {
			//	//TODO 拼接 OAuth2链接，需要AppKey和Secret
			//	url := conf.Config.SysAdmin.Addr + "?redirect_uri="
			//	c.Redirect(http.StatusFound, url)
			//}
		})
	} else if conf.Config.BaseAuth.Enable {
		//开启基本HTTP认证
		//app.Use(gin.BasicAuth(gin.Accounts(conf.Config.BaseAuth.Users)))
	}

	//前端静态文件
	//app.Get("/*any", func(c iris.Context) {
	app.Use(func(c iris.Context) {
		if c.Method() == http.MethodGet {
			//支持前端框架的无“#”路由
			if c.Request().RequestURI == "/" {
				c.Request().URL.Path = "index.html"
			} else if _, err := wwwFiles.FS.Stat(wwwFiles.CTX, c.Request().RequestURI); err != nil {
				c.Request().URL.Path = "index.html"
			}
			//TODO 如果未登录，则跳转SysAdmin OAuth2自动授权页面

			//文件失效期已经在Handler中处理
			wwwFiles.Handler.ServeHTTP(c.ResponseWriter(), c.Request())
		}
	})

	//监听HTTP
	if err := app.Listen(conf.Config.Web.Addr); err != nil {
		log.Fatal("HTTP 服务启动错误", err)
	}
}
