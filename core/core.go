package core

import (
	"errors"
	"git.zgwit.com/zgwit/dtu-admin/db"
	"git.zgwit.com/zgwit/dtu-admin/models"
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
	var cs []models.Tunnel
	err := db.DB("tunnel").All(&cs)
	if err != nil {
		return err
	}

	for _, c := range cs {
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

func StartTunnel(c *models.Tunnel) (Tunnel, error) {
	//log.Println("Start core", c)
	tunnel, err := NewTunnel(c)
	if err != nil {
		return nil, err
	}
	err = tunnel.Open()
	if err != nil {
		return nil, err
	}
	channels.Store(c.ID, tunnel)
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

func GetTunnel(id int) (Tunnel, error) {
	v, ok := channels.Load(id)
	if !ok {
		return nil, errors.New("通道不存在")
	}
	return v.(Tunnel), nil
}

func GetLink(channelId, linkId int) (*Link, error) {
	channel, err := GetTunnel(channelId)
	if err != nil {
		return nil, err
	}
	return channel.GetLink(linkId)
}
