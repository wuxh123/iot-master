package core

import (
	"errors"
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"log"
	"sync"
)

var channels sync.Map

func Tunnels() []Tunnel {
	cs := make([]Tunnel, 0)
	channels.Range(func(key, value interface{}) bool {
		cs = append(cs, value.(Tunnel))
		return true
	})
	return cs
}

func Recovery() error {
	//TODO 改为 加载模型，创建通道
	var cs []models.ModelTunnel
	err := db.Engine.Find(&cs)
	if err != nil {
		return err
	}

	for _, c := range cs {
		//if !c.Disabled {
		_, err = StartTunnel(&c)
		if err != nil {
			log.Println(err)
		}
		//}
	}

	return nil
}

func StartTunnel(c *models.ModelTunnel) (Tunnel, error) {
	//log.Println("Start core", c)
	tunnel, err := NewTunnel(c)
	if err != nil {
		return nil, err
	}
	err = tunnel.Open()
	if err != nil {
		return nil, err
	}
	channels.Store(c.Id, tunnel)
	return tunnel, err
}

func DeleteTunnel(id int) error {
	v, ok := channels.Load(id)
	if !ok {
		return errors.New("通道不存在")
	}
	channels.Delete(id)
	return v.(Tunnel).Close()
}

func GetTunnel(id int64) (Tunnel, error) {
	v, ok := channels.Load(id)
	if !ok {
		return nil, errors.New("通道不存在")
	}
	return v.(Tunnel), nil
}

func GetLink(channelId, linkId int64) (*Link, error) {
	channel, err := GetTunnel(channelId)
	if err != nil {
		return nil, err
	}
	return channel.GetLink(linkId)
}
