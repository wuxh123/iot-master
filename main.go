package main

import (
	"git.zgwit.com/zgwit/iot-admin/conf"
	"git.zgwit.com/zgwit/iot-admin/core"
	"git.zgwit.com/zgwit/iot-admin/flag"
	"git.zgwit.com/zgwit/iot-admin/web"
	"github.com/denisbrodbeck/machineid"
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

	id, err := machineid.ID()
	if err != nil {
		log.Println("获取ID错误：", err)
		return
	}
	log.Println("Machine-Id:", id)

	//加载配置
	err = conf.Load()
	if err != nil {
		log.Println(err)
		return
	}

	//启动总线
	err = core.StartDBus(conf.Config.DBus.Addr)
	if err != nil {
		log.Println(err)
		return
	}

	//恢复之前的链接
	err = core.Recovery()
	if err != nil {
		log.Println(err)
		return
	}

	//启动Web服务
	web.Serve()
}
