package model

import (
	"time"
)

type Link struct {
	Id        int       `json:"id" storm:"id,increment"`
	Name      string    `json:"name"`
	Error     string    `json:"error"`
	Serial    string    `json:"serial" storm:"index"`
	Role      string    `json:"role"`
	Net       string    `json:"net"`
	Addr      string    `json:"addr"`
	ChannelId int     `json:"channel_id"`
	PluginId  int     `json:"plugin_id"` //插件ID
	Online    bool      `json:"online"`
	OnlineAt  time.Time `json:"online_at"`
	CreatedAt time.Time `json:"created_at"`
}
