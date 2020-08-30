package dtu

import (
	"github.com/zgwit/dtu-admin/storage"
	"sync"
)


var channels *sync.Map
var connections *sync.Map

func init() {
	channels = new(sync.Map)
	connections = new(sync.Map)
}

func Channels() *sync.Map {
	return channels
}

func Connections() *sync.Map  {
	return connections
}

func Recovery() error {
	var cs []storage.Channel
	err := storage.ChannelDB().All(&cs)
	if err != nil {
		return err
	}

	for _, c := range cs {
		_, _ = startChannel(&c)
	}

	return nil
}

func startChannel(c *storage.Channel) (*Channel, error) {
	channel := NewChannel(c)
	err := channel.Open()
	if err != nil && channel != nil {
		channel.Error = err.Error()
	}
	channels.Store(c.ID, channel)
	return channel, err
}

func CreateChannel(c *storage.Channel) (*Channel, error)  {
	err := storage.ChannelDB().Save(c)
	if err != nil {
		return nil, err
	}
	return startChannel(c)
}


