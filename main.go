package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/dtu-admin/conf"
	"github.com/zgwit/dtu-admin/flag"
	"github.com/zgwit/dtu-admin/storage"
	"log"
)

func main() {
	//解析参数
	if !flag.Parse() {
		return
	}
	//加载配置
	conf.Load()

	err := storage.Open()
	if err != nil {
		log.Println("数据库错误：", err)
		return
	}

	if !conf.Config.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	//GIN初始化
	app := gin.Default()

	//跨域咨问题
	app.Use(cors.Default())

	//监听HTTP
	if err := app.Run(conf.Config.Http.Addr); err != nil {
		log.Fatal("HTTP 服务启动错误", err)
	}
}
