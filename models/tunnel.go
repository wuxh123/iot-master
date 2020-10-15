package models

import "time"

type Tunnel struct {
	Id      int64 `json:"id"`
	ModelId int64 `json:"model_id"` //模型ID

	Role    string `json:"role"`
	Net     string `json:"net"`
	Addr    string `json:"addr"`
	Timeout int    `json:"timeout"`

	RegisterEnable bool   `json:"register_enable"`
	RegisterRegex  string `json:"register_regex"`
	RegisterMin    int    `json:"register_min"`
	RegisterMax    int    `json:"register_max"`

	HeartBeatEnable   bool   `json:"heart_beat_enable"`
	HeartBeatInterval int    `json:"heart_beat_interval"`
	HeartBeatContent  string `json:"heart_beat_content"`
	HeartBeatIsHex    bool   `json:"heart_beat_is_hex"`

	Disabled bool `json:"disabled"`
	Active   bool `json:"active"` //系统启动要全部置置为false

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

type Link struct {
	Id       int64 `json:"id"`
	TunnelId int64 `json:"tunnel_id"`
	ModelId  int64 `json:"model_id"` //模型ID，默认继承自Tunnel

	Serial string `json:"serial" xorm:"index"`

	Active bool `json:"active"` //系统启动倒置置false

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}
