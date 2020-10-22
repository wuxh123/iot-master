package models

import "time"

type Device struct {
	Id        int64 `json:"id"`
	TunnelId  int64 `json:"tunnel_id"`
	LinkId    int64 `json:"link_id"`
	ProjectId int64 `json:"project_id"`

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
	Id        int64   `json:"id"`
	DeviceId  int64   `json:"device_id"`
	Latitude  float64 `json:"latitude"`  //纬度
	Longitude float64 `json:"longitude"` //经度
	//Altitude  float64   `json:"altitude"`  //高度 单位m
	Created time.Time `json:"created" xorm:"created"`
}
