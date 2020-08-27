package storage

import (
	"github.com/asdine/storm/v3"
	"github.com/zgwit/dtu-admin/conf"
	"os"
	"path/filepath"
)

var (
	channelDB *storm.DB
)

func Open() error {
	var err error
	path := conf.Config.Storage.Path

	err = os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {

		return err
	}

	channelDB, err = storm.Open(filepath.Join(path, "channels.db"))
	if err != nil {
		return err
	}

	return nil
}

func ChannelDB() *storm.DB {
	return channelDB
}
