package model

import "time"

type Plugin struct {
	Id        int64     `json:"id" storm:"increment"`
	Name      string    `json:"name" xorm:"varchar(64)"`
	AppKey    string    `json:"app_key" xorm:"varchar(64)"`
	AppSecret string    `json:"app_secret" xorm:"varchar(64)"`
	Address   string    `json:"address" xorm:"varchar(128)"`
	Path      string    `json:"path" xorm:"varchar(256)"`
	Entry     string    `json:"entry" xorm:"varchar(256)"`
	ExpireAt  time.Time `json:"expire_at"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
}
