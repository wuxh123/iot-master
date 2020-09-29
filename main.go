package main

import (
	"git.zgwit.com/zgwit/iot-admin/flag"
	"git.zgwit.com/zgwit/iot-admin/internal"
	"git.zgwit.com/zgwit/iot-admin/internal/web"
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

	//启动总线
	err = internal.Start()
	if err != nil {
		log.Println("启动失败：", err)
		return
	}

	//启动Web服务
	web.Serve()
}
