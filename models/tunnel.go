package models

import "time"

type Tunnel struct {
	Id            int64 `json:"id"`
	ModelId       int64 `json:"model_id"`
	ModelTunnelId int64 `json:"model_tunnel_id"`

	Active bool `json:"active"` //系统启动倒置置false

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

type Link struct {
	Id            int64 `json:"id"`
	TunnelId      int64 `json:"tunnel_id"`
	ModelId       int64 `json:"model_id"` //冗余，可以通过联合查询
	ModelTunnelId int64 `json:"model_tunnel_id"`

	Serial string `json:"serial" xorm:"index"`

	Active bool `json:"active"` //系统启动倒置置false

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}
