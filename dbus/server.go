package dbus

import (
	"git.zgwit.com/iot/beeq"
	"git.zgwit.com/iot/beeq/packet"
	"log"
)

var hive *beeq.Hive

func Start(addr string) error {
	hive = beeq.NewHive()
	hive.OnConnect(func(connect *packet.Connect, bee *beeq.Bee) bool {
		log.Println(bee.ClientId(), "connect", connect)
		return true
	})
	hive.OnPublish(func(publish *packet.Publish, bee *beeq.Bee) bool {
		log.Println(bee.ClientId(), "publish", publish)
		return true
	})
	return hive.ListenAndServe(addr)
}

func Hive() *beeq.Hive {
	return hive
}
