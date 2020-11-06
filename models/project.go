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

	Created time.Time `json:"created" storm:"created"`
}

type ProjectElement struct {
	ID int `json:"id"`
	//ProjectId     int `json:"project_id"`
	ElementId     int `json:"element_id"`
	ProjectLinkId int `json:"project_link_id"`

	Name  string `json:"name"`
	Alias string `json:"alias"` //别名，用于编程
	Slave uint8  `json:"slave"` //从站号

	Created time.Time `json:"created" storm:"created"`
}

type ProjectVariable struct {
	ID int `json:"id"`
	//ProjectId int `json:"project_id"`
	//ElementId        int `json:"element_id"`
	ProjectElementId int `json:"project_element_id"`

	Variable `storm:"inline"`

	//TODO 添加采样周期

	Created time.Time `json:"created" storm:"created"`
}

type ProjectValidator struct {
	ID         int       `json:"id"`
	ProjectId  int       `json:"project_id"`
	Alert      string    `json:"alert"`
	Expression string    `json:"expression"` //表达式，检测变量名
	Created    time.Time `json:"created" storm:"created"`
}

type ProjectFunction struct {
	ID          int       `json:"id"`
	ProjectId   int       `json:"project_id"`
	Name        string    `json:"name"` //项目功能脚本唯一，供外部调用
	Description string    `json:"description"`
	Script      string    `json:"script"` //javascript
	Created     time.Time `json:"created" storm:"created"`
}

type ProjectJob struct {
	ID         int       `json:"id"`
	ProjectId  int       `json:"project_id"`
	FunctionId string    `json:"function_id"`
	Cron       string    `json:"cron"`
	Created    time.Time `json:"created" storm:"created"`
}

type ProjectStrategy struct {
	ID         int       `json:"id"`
	ProjectId  int       `json:"project_id"`
	FunctionId string    `json:"function_id"`
	Expression string    `json:"expression"` //触发条件 表达式，检测变量名
	Created    time.Time `json:"created" storm:"created"`
}

type Script struct {
	source    string
	variables []string
	script    *otto.Script
}
