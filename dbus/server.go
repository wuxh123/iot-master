package dbus

import (
	"git.zgwit.com/iot/beeq"
	"git.zgwit.com/iot/beeq/packet"
	"log"
	"strconv"
	"strings"
)

var hive *beeq.Hive

func Start(addr string) error {
	hive = beeq.NewHive()
	hive.OnConnect(func(connect *packet.Connect, bee *beeq.Bee) bool {
		//TODO 验证插件 Key Secret
		//TODO 验证浏览器
		//TODO 验证透传
		log.Println(bee.ClientId(), "connect", connect)
		return true
	})
	hive.OnPublish(func(publish *packet.Publish, bee *beeq.Bee) bool {
		log.Println(bee.ClientId(), "publish", publish)
		return true
	})
	hive.Subscribe("/+/send", func(pub *packet.Publish) {
		log.Println(string(pub.Topic()), string(pub.Payload()))

		topics := strings.Split(string(pub.Topic()), "/")
		_, err := strconv.Atoi(topics[1])
		if err != nil {
			log.Println(err)
		}

		//TODO 发送到 link
	})
	return hive.ListenAndServe(addr)
}

func Hive() *beeq.Hive {
	return hive
}
