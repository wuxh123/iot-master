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
	//ElementId int `json:"element_id"`
	Element string `json:"element"` //uuid

	Name  string `json:"name"`
	Alias string `json:"alias"` //别名，用于编程
	Slave uint8  `json:"slave"` //从站号

	Created time.Time `json:"created" storm:"created"`
}

type ProjectVariable struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Alias  string `json:"alias"` //别名，用于编程
	Slave  uint8  `json:"slave"`
	Code   uint8  `json:"code"`   //功能码
	Offset uint16 `json:"offset"` //偏移
	Type   string `json:"type"`
	Unit   string `json:"unit"` //单位

	Scale float32 `json:"scale"` //倍率，比如一般是 整数÷10，得到

	Default  string `json:"default"`
	ReadOnly bool   `json:"read_only"` //只读
}

type ProjectValidator struct {
	ID        int `json:"id"`
	ProjectId int `json:"project_id"`
	//Name       string    `json:"name"`
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
	ID        int       `json:"id"`
	ProjectId int       `json:"project_id"`
	Function  string    `json:"function"`
	Cron      string    `json:"cron"`
	Created   time.Time `json:"created" storm:"created"`
}

type ProjectStrategy struct {
	ID         int       `json:"id"`
	ProjectId  int       `json:"project_id"`
	Expression string    `json:"expression"` //触发条件 表达式，检测变量名
	Function   string    `json:"function"`
	Created    time.Time `json:"created" storm:"created"`
}

type Script struct {
	source    string
	variables []string
	script    *otto.Script
}
