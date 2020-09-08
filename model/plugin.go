package model

import "time"

type Plugin struct {
	Id        int64     `json:"id" storm:"increment"`
	Name      string    `json:"name" xorm:"varchar(64)"`
	AppKey    string    `json:"app_key" xorm:"varchar(64)"`
	AppSecret string    `json:"app_secret" xorm:"varchar(64)"`
	ExpireAt  time.Time `json:"expire_at"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
}
