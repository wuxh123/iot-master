package model

import (
	"time"
)

type Link struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name" xorm:"varchar(64)"`
	Error     string    `json:"error" xorm:"varchar(256)"`
	Serial    string    `json:"serial" xorm:"varchar(128)"`
	Role      string    `json:"role" xorm:"varchar(16) notnull"`
	Net       string    `json:"net" xorm:"varchar(16) notnull"`
	Addr      string    `json:"addr" xorm:"varchar(128) notnull"`
	ChannelId int64     `json:"channel_id"`
	PluginId  int64     `json:"plugin_id"` //插件ID
	Online    bool      `json:"online"`
	OnlineAt  time.Time `json:"online_at"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
}
