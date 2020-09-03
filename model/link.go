package model

import (
	"time"
)

type Link struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name" xorm:"varchar(64)"`
	Serial    string    `json:"serial" xorm:"varchar(128)"`
	Addr      string    `json:"addr" xorm:"varchar(128) notnull"`
	ChannelId int64     `json:"channel_id"`
	PluginId  int64     `json:"plugin_id"` //插件ID
	Online    time.Time `json:"online"`
	Created   time.Time `json:"created" xorm:"created"`
}
