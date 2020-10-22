package open

import (
	"github.com/kataras/iris/v12"
)

func RegisterRoutes(app iris.Party) {

	//跨域咨问题
	//app.Use(cors.Default())

	app.Use(func(ctx iris.Context) {
		//TODO 检查KEY

		//log.Println("open", ctx.FullPath())
	})

	//app.Get("/channels")
}
