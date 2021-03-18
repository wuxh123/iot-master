package model

import (
	"time"
)

type Project struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Disabled bool `json:"disabled"`

	//如果是模板项目，则无效
	Manifest TemplateManifest `json:"manifest" xorm:"json"`

	//模板项目
	TemplateId int64 `json:"template_id"`

	Created time.Time `json:"created" xorm:"created"`
}
