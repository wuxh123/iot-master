package dtu

import (
	"github.com/zgwit/dtu-admin/storage"
	"sync"
)


var channels *sync.Map

func init() {
	channels = new(sync.Map)
}

func Channels() *sync.Map {
	return channels
}

func Recovery() error {
	var cs []storage.Channel
	err := storage.ChannelDB().All(&cs)
	if err != nil {
		return err
	}

	for _, c := range cs {
		channel := NewChannel(c.ID, c.Net, c.Addr, c.IsServer)
		err := channel.Open()
		if err != nil {
			channel.Error = err.Error()
		}

		channels.Store(c.ID, channel)
	}

	return nil
}
