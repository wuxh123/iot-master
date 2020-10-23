package models

import "time"

type Tunnel struct {
	ID      int `json:"id"`
	ModelId int `json:"model_id"` //模型ID

	Name    string `json:"name"`
	Type    string `json:"type"` //tcp-server tcp-client udp-server udp-client serial
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

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

type Link struct {
	ID       int `json:"id"`
	TunnelId int `json:"tunnel_id"`
	ModelId  int `json:"model_id"` //模型ID，默认继承自Tunnel

	Serial string `json:"serial" storm:"index"`

	Active bool `json:"active"` //系统启动倒置置false

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}
