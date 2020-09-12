package db

import (
	"git.zgwit.com/iot/dtu-admin/conf"
	"github.com/zgwit/storm/v3"
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

	path := conf.Config.Data.Path

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
