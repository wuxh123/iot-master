package web

import (
	"github.com/gin-gonic/gin"
	//"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/zgwit/dtu-admin/conf"
	_ "github.com/zgwit/dtu-admin/docs"
	"github.com/zgwit/dtu-admin/web/api"
	"github.com/zgwit/dtu-admin/web/open"
	wwwFiles "github.com/zgwit/dtu-admin/www"
	swaggerFiles "github.com/zgwit/swagger-files"
	"log"
	"net/http"
)


func Serve()  {
	if !conf.Config.Web.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	//GIN初始化
	app := gin.Default()

	api.RegisterRoutes(app.Group("/api"))
	open.RegisterRoutes(app.Group("/open"))

	//加入swagger会增加10MB多体积，使用github.com/zgwit/swagger-files，去除Map文件，可以节省7MB左右
	//Swagger文档，需要先执行swag init生成文档
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//前端静态文件
	//app.GET("/*any", func(c *gin.Context) {
	app.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			//支持前端框架Angular的无“#”路由
			if c.Request.RequestURI == "/" {
				c.Request.URL.Path = "index.html"
			} else if _, err := wwwFiles.FS.Stat(wwwFiles.CTX, c.Request.RequestURI) ; err != nil {
				c.Request.URL.Path = "index.html"
			}
			//文件失效期已经在Handler中处理
			wwwFiles.Handler.ServeHTTP(c.Writer, c.Request)
		}
	})

	//监听HTTP
	if err := app.Run(conf.Config.Web.Addr); err != nil {
		log.Fatal("HTTP 服务启动错误", err)
	}
}
