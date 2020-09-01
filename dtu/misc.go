package dtu

import (
	"errors"
	"github.com/zgwit/dtu-admin/storage"
	"github.com/zgwit/dtu-admin/types"
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
		if !c.Disabled {
			_, _ = startChannel(&c)
		}
	}

	return nil
}

func startChannel(c *types.Channel) (*Channel, error) {
	channel := NewChannel(c)
	err := channel.Open()
	if err != nil && channel != nil {
		channel.Error = err.Error()
	}
	channels.Store(c.ID, channel)
	return channel, err
}

func CreateChannel(c *types.Channel) (*Channel, error)  {
	err := storage.DB("channel").Save(c)
	if err != nil {
		return nil, err
	}
	return startChannel(c)
}

func GetChannel(id int) (*Channel, error) {
	v, ok := channels.Load(id)
	if !ok {
		return nil, errors.New("通道不存在")
	}
	return v.(*Channel), nil
}

func DeleteChannel(id int) error {
	v, ok := channels.Load(id)
	if !ok {
		return errors.New("通道不存在")
	}
	channel := v.(*Channel)
	channel.Close()
	channels.Delete(id)
	return nil
}

