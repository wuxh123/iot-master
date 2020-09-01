package dtu

import (
	"errors"
	"github.com/zgwit/dtu-admin/storage"
	"github.com/zgwit/dtu-admin/types"
	"log"
	"sync"
)


var channels sync.Map
var connections sync.Map

func Channels() []*Channel {
	cs := make([]*Channel, 0)
	channels.Range(func(key, value interface{}) bool {
		cs = append(cs, value.(*Channel))
		return true
	})
	return cs
}


func Recovery() error {
	var cs []types.Channel
	err := storage.DB("channel").All(&cs)
	if err != nil {
		return err
	}

	for _, c := range cs {
		if !c.Net.Disabled {
			_, err = StartChannel(&c)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func StartChannel(c *types.Channel) (*Channel, error) {
	log.Println("start", c)

	channel := NewChannel(c)
	err := channel.Open()
	if err != nil && channel != nil {
		channel.Error = err.Error()
	}
	channels.Store(c.ID, channel)
	return channel, err
}

func GetChannel(id int) (*Channel, error) {
	v, ok := channels.Load(id)
	if !ok {
		return nil, errors.New("通道不存在")
	}
	return v.(*Channel), nil
}
