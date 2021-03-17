package tunnel

import (
	"git.zgwit.com/iot/beeq/packet"
	"iot-master/db"
	"iot-master/dbus"
	"iot-master/model"
	"log"
	"strconv"
	"strings"
)

func Start() error {
	err := Recovery()
	if err != nil {
		return err
	}

	//监听数据
	dbus.Hive().Subscribe("/link/+/+/transfer", func(pub *packet.Publish) {
		//log.Println(string(pub.Topic()), string(pub.Payload()))
		topics := strings.Split(string(pub.Topic()), "/")
		tunnelId, err := strconv.ParseInt(topics[2], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}
		linkId, err := strconv.ParseInt(topics[3], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		//发送到 link
		link, err := GetLink(tunnelId, linkId)
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

	return nil
}

func Recovery() error {
	var ts []model.Tunnel
	err := db.Engine.Find(&ts)
	if err != nil {
		return err
	}

	for _, c := range ts {
		if c.Disabled {
			continue
		}
		_, err = StartTunnel(&c)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
