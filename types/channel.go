package types

import "time"

type RegisterConf struct {
	Enable bool
	Length int
	Regex  string
}

type HeartBeatConf struct {
	Enable  bool
	Content []byte
}

type Channel struct {
	ID   int      `storm:"increment" json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`

	Serial string `storm:"index" json:"serial"`

	Net      string `json:"net"`
	Addr     string `json:"addr"`
	IsServer bool   `json:"is_server"`
	Timeout  int    `json:"timeout"` //TODO 改为秒

	Register struct {
		Enable bool   `json:"enable"`
		Length int    `json:"length"`
		Regex  string `json:"regex"`
	} `json:"register"`

	HeartBeat struct {
		Enable   bool   `json:"enable"`
		Interval int    `json:"interval"` //TODO 改为秒
		Content  []byte `json:"content"`
	} `json:"heart_beat"`

	Disabled bool      `json:"disabled"`
	Created  time.Time `json:"created"`
	Creator  int       `json:"creator"`
}
