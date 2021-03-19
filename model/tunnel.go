package model

import "time"

type Register struct {
	Enable bool   `json:"enable"`
	Regex  string `json:"regex"`
	Min    int    `json:"min"`
	Max    int    `json:"max"`
	//添加每次注册包检测，暂不需要（UDP可能会用）
}

type HeartBeat struct {
	Enable   bool   `json:"enable"`
	Interval int    `json:"interval"`
	Content  string `json:"content"`
	IsHex    bool   `json:"is_hex"`
}

type Tunnel struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"` //tcp-server tcp-client udp-server udp-client serial
	Addr    string `json:"addr"`
	Timeout int    `json:"timeout"`

	//注册包
	Register Register `json:"register" xorm:"json"`

	//心跳包
	HeartBeat HeartBeat `json:"heart_beat" xorm:"json"`

	//模板ID，根据模板ID自动创建项目
	TemplateId int64 `json:"template_id"`

	Disabled bool `json:"disabled"`

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

type Link struct {
	Id       int64 `json:"id"`
	TunnelId int64 `json:"tunnel_id"`

	ProjectId  int64  `json:"project_id"`  //项目ID
	ProjectKey string `json:"project_key"` //项目中 链接KEY

	Serial string `json:"serial" xorm:"index"`
	Addr   string `json:"addr"`

	Active  bool      `json:"active"`
	Online  time.Time `json:"online"`
	Created time.Time `json:"created" xorm:"created"`
}
