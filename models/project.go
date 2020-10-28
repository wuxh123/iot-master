package models

import (
	"time"
	"github.com/robertkrimen/otto"
)

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

	Links      []ProjectLink      `json:"links"`
	Jobs       []ProjectJob       `json:"jobs"`
	Strategies []ProjectStrategy  `json:"strategies"`
	Validators []ProjectValidator `json:"validators"`

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

type ProjectLink struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`

	Elements []ProjectElement `json:"elements"`
}

type ProjectElement struct {
	ElementId int `json:"element_id"`

	Name  string `json:"name"`
	Slave uint8  `json:"slave"` //从站号

	Variables []ProjectVariable `json:"variables"`
}

type ProjectVariable struct {
	ElementVariable `storm:"inline"`

	Name    string  `json:"name"`
	Alias   string  `json:"alias"`   //别名，用于编程
	Correct float32 `json:"correct"` //校准
}

type ProjectValidator struct {
	Message string `json:"message"`
	Script  string `json:"script"`
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

type Script struct {
	source    string
	variables []string
	script    *otto.Script
}
