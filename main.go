package main

import (
	"github.com/denisbrodbeck/machineid"
	"git.zgwit.com/iot/dtu-admin/conf"
	"git.zgwit.com/iot/dtu-admin/dbus"
	"git.zgwit.com/iot/dtu-admin/dtu"
	"git.zgwit.com/iot/dtu-admin/flag"
	"git.zgwit.com/iot/dtu-admin/web"
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

	id, err := machineid.ID()
	if err != nil {
		log.Println("获取ID错误：", err)
		return
	}
	log.Println("Machine Id:", id)

	//
	//err = db.Open()
	//if err != nil {
	//	log.Println("数据库错误：", err)
	//	return
	//}

	//启动总线
	err = dbus.Start(conf.Config.DBus.Addr)
	if err != nil {
		log.Println("总线启动失败：", err)
		return
	}

	//恢复之前的链接
	err = dtu.Recovery()
	if err != nil {
		log.Println("恢复链接：", err)
		return
	}


	//启动Web服务
	web.Serve()
}
