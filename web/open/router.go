package open

import (
	"github.com/gin-gonic/gin"
	"log"
)

func RegisterRoutes(app *gin.RouterGroup) {
	app.Use(func(ctx *gin.Context) {
		//TODO 检查KEY

		log.Println("open", ctx.FullPath())
	})

	app.GET("/channels")
}
