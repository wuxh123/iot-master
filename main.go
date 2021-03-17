package main

import (
	"iot-master/args"
	"iot-master/conf"
	"iot-master/db"
	"iot-master/dbus"
	_ "iot-master/protocol/modbus" //默认支持Modbus协议
	"iot-master/tunnel"
	"iot-master/web"
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
	args.Parse()

	//加载配置
	err := conf.Load()
	if err != nil {
		log.Println(err)
		return
	}

	err = db.Open()
	if err != nil {
		log.Println("数据库错误：", err)
		return
	}

	//启动总线
	err = dbus.Start(conf.Config.DBus.Addr)
	if err != nil {
		log.Println(err)
		return
	}

	//恢复之前的链接
	err = tunnel.Recovery()
	if err != nil {
		log.Println(err)
		return
	}

	//启动Web服务
	web.Serve()
}
