package open

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *gin.RouterGroup) {

	//跨域咨问题
	app.Use(cors.Default())

	app.Use(func(ctx *gin.Context) {
		//TODO 检查KEY

		//log.Println("open", ctx.FullPath())
	})

	app.GET("/channels")
}
