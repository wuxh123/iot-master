package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/zgwit/dtu-admin/conf"
	_ "github.com/zgwit/dtu-admin/docs"
	"github.com/zgwit/dtu-admin/web/api"
	"github.com/zgwit/dtu-admin/web/open"
	"log"
	"net/http"
)

func Serve()  {
	if !conf.Config.Web.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	//GIN初始化
	app := gin.Default()

	//跨域咨问题
	app.Use(cors.Default())

	//前端静态文件
	app.Use(static.Serve("/", static.LocalFile("./www/", false)))


	api.RegisterRoutes(app.Group("/api"))
	open.RegisterRoutes(app.Group("/open"))


	//Swagger文档，需要先执行swag init生成文档
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//支持前端框架Angular的无“#”路由
	//app.GET("/*any", func(c *gin.Context) {
	app.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			http.ServeFile(c.Writer, c.Request, "./www/index.html")
		}
	})

	//监听HTTP
	if err := app.Run(conf.Config.Web.Addr); err != nil {
		log.Fatal("HTTP 服务启动错误", err)
	}
}
