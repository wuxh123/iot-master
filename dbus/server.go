package dbus

import (
	"git.zgwit.com/iot/beeq"
)

var hive *beeq.Hive

func Start(addr string) error {
	hive = beeq.NewHive()
	return hive.ListenAndServe(addr)
}

func Hive() *beeq.Hive  {
	return hive
}