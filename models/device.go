package models

import "time"

type Device struct {
	Id       int64 `json:"id"`
	TunnelId int64 `json:"tunnel_id"`
	LinkId   int64 `json:"link_id"`
	ModelId  int64 `json:"model_id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Serial      string `json:"serial"`

	//定位，或手动选择
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}
