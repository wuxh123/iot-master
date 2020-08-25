package database

import (
	"github.com/zgwit/dtu-admin/conf"
	//"github.com/zgwit/dtu-admin/model"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

func OpenMySQL() error {

	if Engine != nil {
		return nil
	}

	cfg := conf.Config.Database
	var err error
	Engine, err = xorm.NewEngine(cfg.Type, cfg.Url)
	if err != nil {
		return err
	}
	Engine.ShowSQL(cfg.ShowSQL)



	return nil
}

func SyncMySQL() error {
	//同步表
	return Engine.Sync2()
}