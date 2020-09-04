package model

import "time"

type Channel struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" xorm:"varchar(64)"`
	Error string `json:"error" xorm:"varchar(256)"`

	Disabled bool   `json:"disabled" xorm:"default 0"` //此处 禁用 直接放到顶级，Update无效
	Role     string `json:"role" xorm:"varchar(16) notnull"`
	Net      string `json:"net" xorm:"varchar(16) notnull"`
	Addr     string `json:"addr" xorm:"varchar(128) notnull"`
	Timeout  int    `json:"timeout"`

	RegisterEnable bool   `json:"register_enable" xorm:"default 0"`
	RegisterRegex  string `json:"register_regex" xorm:"varchar(128)"`
	RegisterMin    int    `json:"register_min"`
	RegisterMax    int    `json:"register_max"`

	HeartBeatEnable   bool   `json:"heart_beat_enable" xorm:"default 0"`
	HeartBeatInterval int    `json:"heart_beat_interval"`
	HeartBeatContent  string `json:"heart_beat_content" xorm:"varchar(256)"`
	HeartBeatIsHex    bool   `json:"heart_beat_is_hex" xorm:"default 0"`

	PluginId int64 `json:"plugin_id"`

	//Creator int       `json:"creator"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
}