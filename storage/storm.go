package storage

import (
	"github.com/asdine/storm/v3"
	"github.com/zgwit/dtu-admin/conf"
	"log"
	"os"
	"path/filepath"
	"sync"
)


var databases sync.Map

func DB(name string) *storm.DB {
	if v, ok := databases.Load(name); ok {
		return v.(*storm.DB)
	}

	path := conf.Config.Storage.Path

	err := os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	db, err := storm.Open(filepath.Join(path, name+".db"))
	if err != nil {
		log.Fatal(err)
	}

	databases.Store(name, db)
	return db
}
