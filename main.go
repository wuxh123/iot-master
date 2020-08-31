package main

import (
	"github.com/zgwit/dtu-admin/conf"
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/flag"
	"github.com/zgwit/dtu-admin/web"
	"log"
)

// @title DTU manager API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.zgwit.com/support
// @contact.email jason@zgwit.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1
// @BasePath /open
func main() {
	//解析参数
	if !flag.Parse() {
		return
	}
	//加载配置
	conf.Load()

	//恢复之前的链接
	err := dtu.Recovery()
	if err != nil {
		log.Println("恢复链接：", err)
		return
	}

	web.Serve()
}
