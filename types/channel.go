package types

import "time"

type NetConf struct {
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
	Content  []byte `json:"content"`
}

type Channel struct {
	ID   int      `storm:"increment" json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`

	Net       NetConf       `json:"net"`
	Register  RegisterConf  `json:"register"`
	HeartBeat HeartBeatConf `json:"heart_beat"`

	Disabled bool      `json:"disabled"`
	Created  time.Time `json:"created"`
	Creator  int       `json:"creator"`
}
