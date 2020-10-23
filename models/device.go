package models

import "time"

type Device struct {
	ID        int `json:"id"`
	TunnelId  int `json:"tunnel_id"`
	LinkId    int `json:"link_id"`
	ProjectId int `json:"project_id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Serial      string `json:"serial"`

	//定位，或手动选择
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

//默认WGS84标准，GCJ02、BD09都需要转换
type Location struct {
	ID        int   `json:"id"`
	DeviceId  int   `json:"device_id"`
	Latitude  float64 `json:"latitude"`  //纬度
	Longitude float64 `json:"longitude"` //经度
	//Altitude  float64   `json:"altitude"`  //高度 单位m
	Created time.Time `json:"created" xorm:"created"`
}
