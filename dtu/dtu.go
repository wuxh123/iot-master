package dtu

import (
	"errors"
	"git.zgwit.com/iot/dtu-admin/db"
	"git.zgwit.com/iot/dtu-admin/model"
	"log"
	"sync"
)


var channels sync.Map

func Channels() []Channel {
	cs := make([]Channel, 0)
	channels.Range(func(key, value interface{}) bool {
		cs = append(cs, value.(Channel))
		return true
	})
	return cs
}


func Recovery() error {
	var cs []model.Channel
	err := db.DB("channel").All(&cs)
	if err != nil {
		return err
	}

	for _, c := range cs {
		if !c.Disabled {
			_, err = StartChannel(&c)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func StartChannel(c *model.Channel) (Channel, error) {
	//log.Println("Start channel", c)
	channel, err := NewChannel(c)
	if err != nil {
		return nil, err
	}
	err = channel.Open()
	if err != nil {
		return nil, err
	}
	channels.Store(c.Id, channel)
	return channel, err
}

func DeleteChannel(id int) error  {
	v, ok := channels.Load(id)
	if !ok {
		return errors.New("通道不存在")
	}
	channels.Delete(id)
	return v.(Channel).Close()
}

func GetChannel(id int) (Channel, error) {
	v, ok := channels.Load(id)
	if !ok {
		return nil, errors.New("通道不存在")
	}
	return v.(Channel), nil
}

func GetLink(channelId, linkId int) (*Link, error)  {
	channel, err := GetChannel(channelId)
	if err != nil {
		return nil, err
	}
	return channel.GetLink(linkId)
}
