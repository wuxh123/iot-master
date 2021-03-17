package model

import "time"

type Tunnel struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"` //tcp-server tcp-client udp-server udp-client serial
	Addr    string `json:"addr"`
	Timeout int    `json:"timeout"`

	//注册包
	RegisterEnable bool   `json:"register_enable"`
	RegisterRegex  string `json:"register_regex"`
	RegisterMin    int    `json:"register_min"`
	RegisterMax    int    `json:"register_max"`

	//心跳包
	HeartBeatEnable   bool   `json:"heart_beat_enable"`
	HeartBeatInterval int    `json:"heart_beat_interval"`
	HeartBeatContent  string `json:"heart_beat_content"`
	HeartBeatIsHex    bool   `json:"heart_beat_is_hex"`

	//模板ID，根据模板ID自动创建项目
	TemplateId int64 `json:"template_id"`

	Disabled bool `json:"disabled"`

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

type Link struct {
	Id       int64 `json:"id"`
	TunnelId int64 `json:"tunnel_id"`

	ProjectId int64 `json:"project_id"` //项目ID

	Serial string `json:"serial" storm:"index"`
	Addr   string `json:"addr"`

	Active  bool      `json:"active"`
	Online  time.Time `json:"online"`
	Created time.Time `json:"created" xorm:"created"`
}
