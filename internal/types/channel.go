package types

import "time"

type Statistic struct {
	Rx int `json:"rx"`
	Tx int `json:"tx"`
}

type Channel struct {
	Id    int    `json:"id" storm:"id,increment"`
	Name  string `json:"name"`
	Error string `json:"error"`

	Disabled bool   `json:"disabled"`
	Role     string `json:"role"`
	Net      string `json:"net"`
	Addr     string `json:"addr"`
	Timeout  int    `json:"timeout"`

	RegisterEnable bool   `json:"register_enable"`
	RegisterRegex  string `json:"register_regex"`
	RegisterMin    int    `json:"register_min"`
	RegisterMax    int    `json:"register_max"`

	HeartBeatEnable   bool   `json:"heart_beat_enable"`
	HeartBeatInterval int    `json:"heart_beat_interval"`
	HeartBeatContent  string `json:"heart_beat_content"`
	HeartBeatIsHex    bool   `json:"heart_beat_is_hex"`

	//PluginId int `json:"plugin_id"`
	//TODO 默认模型

	CreatedAt time.Time `json:"created_at"`
}

type ChannelExt struct {
	Channel `storm:"inline"`
	Statistic

	Links  int  `json:"links"`
	Online bool `json:"online"`
}

type Link struct {
	Id        int       `json:"id" storm:"id,increment"`
	Name      string    `json:"name"`
	Error     string    `json:"error"`
	Serial    string    `json:"serial" storm:"index"`
	Net       string    `json:"net"`
	Addr      string    `json:"addr"`
	ChannelId int       `json:"channel_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkExt struct {
	Link `storm:"inline"`
	Statistic

	Online  bool `json:"online"`
}
