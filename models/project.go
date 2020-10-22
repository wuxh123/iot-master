package models

import "time"

type Project struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"` //模板ID
	UUID        string `json:"uuid"`   //唯一码，自动生成

	Version  string    `json:"version"`
	Disabled bool      `json:"disabled"`
	Created  time.Time `json:"created" xorm:"created"`
	Updated  time.Time `json:"updated" xorm:"updated"`
	Deployed time.Time `json:"deployed"` //如果 deployed < updated，说明有更新，提示重新部署
}

type ProjectElement struct {
	Id        int64 `json:"id"`
	ProjectId int64 `json:"project_id"`
	ElementId int64 `json:"element_id"`
	TunnelId  int64 `json:"tunnel_id"`

	Name string `json:"name"`

	Slave uint8  `json:"slave"` //从站号
	Alias string `json:"alias"` //别名，用于编程

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

//
//type ProjectElementVariable struct {
//	Id                int64 `json:"id"` //TODO 去掉ID，用双主键
//	ProjectElementId  int64 `json:"project_element_id"`
//	ElementVariableId int64 `json:"element_variable_id"`
//
//	Created time.Time `json:"created" xorm:"created"`
//}

type ProjectJob struct {
	Id        int64  `json:"id"`
	ProjectId int64  `json:"project_id"`
	Name      string `json:"name"`
	Cron      string `json:"cron"`
	Script    string `json:"script"` //javascript

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

type ProjectStrategy struct {
	Id        int64  `json:"id"`
	ProjectId int64  `json:"project_id"`
	Name      string `json:"name"`
	Script    string `json:"script"` //javascript

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}
