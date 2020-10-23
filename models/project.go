package models

import "time"

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"` //模板ID
	UUId        string `json:"uuId"`   //唯一码，自动生成

	Version  string `json:"version"`
	Disabled bool   `json:"disabled"`

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

type ProjectElement struct {
	ID        int `json:"id"`
	ProjectId int `json:"project_id"`
	ElementId int `json:"element_id"`
	TunnelId  int `json:"tunnel_id"`

	Name string `json:"name"`

	Slave uint8  `json:"slave"` //从站号
	Alias string `json:"alias"` //别名，用于编程

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

//
//type ProjectElementVariable struct {
//	ID                int `json:"id"` //TODO 去掉ID，用双主键
//	ProjectElementId  int `json:"project_element_id"`
//	ElementVariableId int `json:"element_variable_id"`
//
//	Created time.Time `json:"created" storm:"created"`
//}

type ProjectJob struct {
	ID        int    `json:"id"`
	ProjectId int    `json:"project_id"`
	Name      string `json:"name"`
	Cron      string `json:"cron"`
	Script    string `json:"script"` //javascript

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

type ProjectStrategy struct {
	ID        int    `json:"id"`
	ProjectId int    `json:"project_id"`
	Name      string `json:"name"`
	Script    string `json:"script"` //javascript

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}
