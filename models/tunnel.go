package models

import "time"

type Tunnel struct {
	ID      int    `json:"id"`
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

	Active bool `json:"active"`

	TemplateId int `json:"template_id"` //模板ID

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

type Link struct {
	ID        int `json:"id"`
	TunnelId  int `json:"tunnel_id"`

	ProjectId int `json:"project_id"` //项目ID

	Serial string `json:"serial" storm:"index"`
	Addr   string `json:"addr"`

	Active  bool      `json:"active"`
	Online  time.Time `json:"online"`
	Created time.Time `json:"created" storm:"created"`
}
