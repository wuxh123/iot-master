package model

import "time"

type Device struct {
	Id        int64 `json:"id"`
	ProjectId int64 `json:"project_id"`
	Link      uint8 `json:"link"`

	Slave uint8 `json:"slave"` //从站号

	Created time.Time `json:"created" xorm:"created"`
}
