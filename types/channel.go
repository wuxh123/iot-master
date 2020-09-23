package types

import "time"

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

	CreatedAt time.Time `json:"created_at"`
}
