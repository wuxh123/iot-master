package core

import (
	"git.zgwit.com/iot/beeq"
	"git.zgwit.com/iot/beeq/packet"
	"mydtu/db"
	"mydtu/model"
	"github.com/zgwit/storm/v3"
	"log"
	"strconv"
	"strings"
)

var hive *beeq.Hive

func StartDBus(addr string) error {
	hive = beeq.NewHive()
	hive.OnConnect(func(connect *packet.Connect, bee *beeq.Bee) bool {
		// 验证插件 Key Secret
		var plugin model.Plugin
		err := db.DB("plugin").One("Key", connect.UserName(), &plugin)
		if err == storm.ErrNotFound {
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

	hive.Subscribe("/link/+/+/transfer", func(pub *packet.Publish) {
		//log.Println(string(pub.Topic()), string(pub.Payload()))
		topics := strings.Split(string(pub.Topic()), "/")
		channelId, err := strconv.Atoi(topics[2])
		if err != nil {
			log.Println(err)
			return
		}
		linkId, err := strconv.Atoi(topics[3])
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
		err = link.Write(pub.Payload())
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
