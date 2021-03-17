package dbus

import (
	"git.zgwit.com/iot/beeq"
	"git.zgwit.com/iot/beeq/packet"
	"iot-master/db"
	"iot-master/model"
	"log"
)

var hive *beeq.Hive

func Start(addr string) error {
	hive = beeq.NewHive()
	hive.OnConnect(func(connect *packet.Connect, bee *beeq.Bee) bool {
		// 验证插件 Key Secret
		var plugin model.Plugin
		has, err := db.Engine.Where("key=?", connect.UserName()).Get(&plugin)
		if !has {
			if plugin.Secret == string(connect.Password()) {
				return true
			} else {
				return false
			}
		} else if err != nil {
			log.Println(err)
			return false
		}

		//TODO 验证浏览器

		log.Println(bee.ClientId(), "connect", connect)
		return true
	})

	//hive.OnPublish(func(publish *packet.Publish, bee *beeq.Bee) bool {
	//	log.Println(bee.ClientId(), "publish", publish)
	//	return true
	//})

	return hive.ListenAndServe(addr)
}

func Hive() *beeq.Hive {
	return hive
}
