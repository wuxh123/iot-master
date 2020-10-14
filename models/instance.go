package models

import "time"

//TODO 修改为设备
type Instance struct {
	Id      int64 `json:"id"`
	ModelId int64 `json:"model_id"`

	Created time.Time `json:"created" xorm:"created"`
}
