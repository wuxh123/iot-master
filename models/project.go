package models

import (
	"github.com/robertkrimen/otto"
	"time"
)

type Project struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid" storm:"unique"` //唯一码，自动生成
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Disabled bool `json:"disabled"`

	Created time.Time `json:"created" storm:"created"`
}

type ProjectLink struct {
	ID        int    `json:"id"`
	ProjectId int    `json:"project_id"`
	LinkId    int    `json:"link_id"`
	Name      string `json:"name"`
	Protocol  string `json:"protocol"`

	//Elements []ProjectElement `json:"elements"`
	Created time.Time `json:"created" storm:"created"`
}

type ProjectElement struct {
	ID        int `json:"id"`
	ProjectId int `json:"project_id"`
	ElementId int `json:"element_id"`
	//Element string `json:"element"` //uuid

	Name  string `json:"name"`
	Alias string `json:"alias"` //项目元件唯一
	Slave uint8  `json:"slave"` //从站号

	//采样周期，使用定时器。如果空闲，则读取， 如果忙，则排队， 如果已经排队，则跳过
	Sampling int `json:"sampling"` //采样周期 ms

	Created time.Time `json:"created" storm:"created"`
}

type ProjectValidator struct {
	ID        int       `json:"id"`
	ProjectId int       `json:"project_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Script    string    `json:"script"`
	Created   time.Time `json:"created" storm:"created"`
}

type ProjectJob struct {
	ID        int       `json:"id"`
	ProjectId int       `json:"project_id"`
	Name      string    `json:"name"`
	Cron      string    `json:"cron"`
	Script    string    `json:"script"` //javascript
	Created   time.Time `json:"created" storm:"created"`
}

type ProjectStrategy struct {
	ID        int    `json:"id"`
	ProjectId int    `json:"project_id"`
	Name      string `json:"name"`
	Trigger   string `json:"trigger"` //触发条件，当条件满足时，执行Script
	//Triggers  string    `json:"triggers"` //触发变量
	Script  string    `json:"script"` //javascript
	Created time.Time `json:"created" storm:"created"`
}

type ProjectFunction struct {
	ID        int       `json:"id"`
	ProjectId int       `json:"project_id"`
	Name      string    `json:"name"`
	Alias     string    `json:"alias"`  //项目功能脚本唯一，供外部调用
	Script    string    `json:"script"` //javascript
	Created   time.Time `json:"created" storm:"created"`
}

type Script struct {
	source    string
	variables []string
	script    *otto.Script
}
