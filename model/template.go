package model

import (
	"time"
)

type Template struct {
	Id     int64  `json:"id"`
	Origin string `json:"origin"` //平台UUID，自动生成

	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Manifest TemplateManifest `json:"manifest" xorm:"json"`

	Created time.Time `json:"created" xorm:"created"`
}

type TemplateManifest struct {
	Links   map[string]TemplateLink   `json:"links"`
	Devices map[string]TemplateDevice `json:"devices"`

	Functions  map[string]TemplateFunction `json:"functions"`
	Strategies map[string]TemplateStrategy `json:"strategies"`

	Validators []TemplateValidator `json:"validators"`
}

type TemplateLink struct {
	Description string `json:"description"`
	//协议
	Protocol        string `json:"protocol"`
	ProtocolOptions string `json:"protocol_options"`
}

type TemplateDevice struct {
	Description string `json:"description"`

	Element string `json:"element"` //uuid

	Link  uint8 `json:"link"`  //链接编号 0 1 2 3
	Slave uint8 `json:"slave"` //从站号

	DeviceId int64 `json:"device_id"` //TODO 模板中不需要

	Variables  map[string]TemplateVariable `json:"variables"`
	Validators []TemplateValidator         `json:"validators"`
}

type TemplateVariable struct {
	Element string `json:"element"` //uuid

	Variable

	//TODO 添加采样周期
}

type TemplateValidator struct {
	Alert      string   `json:"alert"`
	Watch      []string `json:"watch"` //监听变量（前端最好能检索生成）
	Expression string   `json:"expression"` //表达式，检测变量名
}

type TemplateFunction struct {
	Description string            `json:"description"`
	Operators   map[string]string `json:"operators"`
}

type TemplateStrategy struct {
	Description string            `json:"description"`
	Cron        string            `json:"cron"`
	Expression  string            `json:"expression"` //触发条件 表达式，检测变量名
	Operators   map[string]string `json:"operators"`
}
