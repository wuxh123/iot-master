package models

import (
	"github.com/robertkrimen/otto"
	"time"
)

type Project struct {
	ProjectTemplate `storm:"inline"`

	//Disabled   bool  `json:"disabled"`
	TemplateId int   `json:"template_id"`
	LinkBinds  []int `json:"link_binds"`
}

type ProjectTemplate struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid" storm:"unique"` //唯一码，自动生成
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Links      []ProjectLink      `json:"links"`
	Jobs       []ProjectJob       `json:"jobs"`       //定时任务
	Strategies []ProjectStrategy  `json:"strategies"` //策略
	Functions  []ProjectFunction  `json:"functions"`  //功能脚本，比如：批量开启/关闭，修改模式
	Validators []ProjectValidator `json:"validators"` //报警检查

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

type ProjectLink struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`

	Elements []ProjectElement `json:"elements"`
}

type ProjectElement struct {
	Element string `json:"element"` //uuid

	Name  string `json:"name"`
	Alias string `json:"alias"` //项目元件唯一
	Slave uint8  `json:"slave"` //从站号

	//采样周期，使用定时器。如果空闲，则读取， 如果忙，则排队， 如果已经排队，则跳过
	Sampling int `json:"sampling"` //采样周期 ms
}

type ProjectValidator struct {
	Title   string `json:"title"`
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

type ProjectFunction struct {
	Name   string `json:"name"`
	Alias  string `json:"alias"`  //项目功能脚本唯一
	Script string `json:"script"` //javascript
}

type Script struct {
	source    string
	variables []string
	script    *otto.Script
}
