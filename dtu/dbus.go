package dtu

import (
	"git.zgwit.com/iot/beeq"
	"git.zgwit.com/iot/beeq/packet"
	"log"
	"strconv"
	"strings"
)

var hive *beeq.Hive

func StartDBus(addr string) error {
	hive = beeq.NewHive()
	hive.OnConnect(func(connect *packet.Connect, bee *beeq.Bee) bool {
		//TODO 验证插件 Key Secret
		//TODO 验证浏览器
		//TODO 验证透传
		log.Println(bee.ClientId(), "connect", connect)
		return true
	})
	//hive.OnPublish(func(publish *packet.Publish, bee *beeq.Bee) bool {
	//	log.Println(bee.ClientId(), "publish", publish)
	//	return true
	//})
	hive.Subscribe("/+/+/transfer", func(pub *packet.Publish) {
		log.Println(string(pub.Topic()), string(pub.Payload()))

		topics := strings.Split(string(pub.Topic()), "/")
		channelId, err := strconv.Atoi(topics[1])
		if err != nil {
			log.Println(err)
			return
		}
		linkId, err := strconv.Atoi(topics[2])
		if err != nil {
			log.Println(err)
			return
		}

		//发送到 link
		link, err := GetLink(channelId, linkId)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = link.Send(pub.Payload())
		if err != nil {
			log.Println(err)
			return
		}
	})
	return hive.ListenAndServe(addr)
}

func Hive() *beeq.Hive {
	return hive
}
