package types

import "time"

type NetConf struct {
	Disabled bool   `json:"disabled"` //此处 禁用 直接放到顶级，Update无效
	Type     string `json:"type"`
	Addr     string `json:"addr"`
	IsServer bool   `json:"is_server"`
	Timeout  int    `json:"timeout"` //TODO 改为秒
}

type RegisterConf struct {
	Enable bool   `json:"enable"`
	Length int    `json:"length"`
	Regex  string `json:"regex"`
}

type HeartBeatConf struct {
	Enable   bool   `json:"enable"`
	Interval int    `json:"interval"` //TODO 改为秒
	Content  string `json:"content"`
	IsHex    bool   `json:"is_hex"`
}

type Channel struct {
	ID   int      `storm:"increment" json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`

	Net       NetConf       `json:"net"`
	Register  RegisterConf  `json:"register"`
	HeartBeat HeartBeatConf `json:"heart_beat"`

	Plugin int `json:"plugin"` //插件ID TODO，将子连接所有内容转发至该插件

	Created time.Time `json:"created"`
	Creator int       `json:"creator"`
}
