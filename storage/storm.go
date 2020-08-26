package storage

import (
	"github.com/asdine/storm/v3"
	"github.com/zgwit/dtu-admin/conf"
	"os"
	"path/filepath"
)

var (
	Links *storm.DB
)

func Open() error {
	var err error
	path := conf.Config.Storage.Path

	err = os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {

		return err
	}

	Links, err = storm.Open(filepath.Join(path, "links.db"))
	if err != nil {
		return err
	}


	return nil
}