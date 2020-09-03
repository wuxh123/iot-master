package model

import "time"

type Channel struct {
	Id   int64  `json:"id"`
	Name string `json:"name" xorm:"varchar(64)"`
	//Tags string `json:"tags" xorm:"varchar(256)"`

	Disabled bool   `json:"disabled" xorm:"default 0"` //此处 禁用 直接放到顶级，Update无效
	Type     string `json:"type" xorm:"varchar(16) notnull"`
	Addr     string `json:"addr" xorm:"varchar(128) notnull"`
	IsServer bool   `json:"is_server" xorm:"default 0"`
	Timeout  int    `json:"timeout"` //TODO 改为秒

	RegisterEnable bool   `json:"register_enable" xorm:"default 0"`
	RegisterRegex  string `json:"register_regex" xorm:"varchar(128)"`

	HeartBeatEnable   bool   `json:"heart_beat_enable" xorm:"default 0"`
	HeartBeatInterval int    `json:"heart_beat_interval"` //TODO 改为秒
	HeartBeatContent  string `json:"heart_beat_content" xorm:"varchar(256)"`
	HeartBeatIsHex    bool   `json:"heart_beat_is_hex" xorm:"default 0"`

	PluginId int64 `json:"plugin_id"`

	//Creator int       `json:"creator"`
	Created time.Time `json:"created" xorm:"created"`
}
