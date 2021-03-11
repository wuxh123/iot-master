package db

import (
	//_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mydtu/conf"
	"mydtu/model"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

func Open() error {

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

	//同步表
	err = Engine.Sync2(model.User{}, model.User{})
	if err != nil {
		return err
	}

	initial()

	return nil
}

func initial() {
	var u model.User
	has, err := Engine.Where("username=?", "admin").Exist(&u)
	if err != nil {
		log.Println("query user", err)
		return
	}
	if !has {
		u.Username = "admin"
		u.Password = "123456"
		u.Name = "管理员"
		_, _ = Engine.Insert(u)
	}
}
