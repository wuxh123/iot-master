package models

import "time"

type Project struct {
	ProjectTemplate `storm:"inline"`

	TemplateId int   `json:"template_id"`
	LinkBinds  []int `json:"link_binds"`
}

type ProjectTemplate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//UUID        string `json:"uuid"` //唯一码，自动生成

	Version  string `json:"version"`
	Disabled bool   `json:"disabled"`

	Links      []ProjectLink     `json:"links"`
	Elements   []ProjectElement  `json:"elements"`
	Jobs       []ProjectJob      `json:"jobs"`
	Strategies []ProjectStrategy `json:"strategies"`

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

type ProjectLink struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
}

type ProjectElement struct {
	ElementId int `json:"element_id"`

	Link int `json:"link"` //链接号：0,1,2,3

	Slave uint8  `json:"slave"` //从站号
	Name  string `json:"name"`
	Alias string `json:"alias"` //别名，用于编程

	Variables []ProjectElementVariable `json:"variables"`
}

type ProjectElementVariable struct {
	ElementVariable `storm:"inline"`

	Name    string  `json:"name"`
	Alias   string  `json:"alias"`   //别名，用于编程
	Correct float32 `json:"correct"` //校准
}

type ProjectJob struct {
	Name   string `json:"name"`
	Cron   string `json:"cron"`
	Script string `json:"script"` //javascript
}

type ProjectStrategy struct {
	Name   string `json:"name"`
	Script string `json:"script"` //javascript

}
