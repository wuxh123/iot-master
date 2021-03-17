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

	//协议
	Protocol        string `json:"protocol"`
	ProtocolOptions string `json:"protocol_options"`

	//变量等
	Variables  []TemplateVariable  `json:"variables"`
	Validators []TemplateValidator `json:"validators"`
	Functions  []TemplateFunction  `json:"functions"`
	Strategies []TemplateStrategy  `json:"strategies"`
}

type TemplateVariable struct {
	Element string `json:"element"` //uuid
	Slave   uint8  `json:"slave"`   //从站号

	Variable

	//TODO 添加采样周期
}

type TemplateValidator struct {
	Alert      string `json:"alert"`
	Expression string `json:"expression"` //表达式，检测变量名
}

type TemplateFunction struct {
	Name        string                 `json:"name"` //项目功能脚本唯一，供外部调用
	Description string                 `json:"description"`
	Script      string                 `json:"script"` //javascript
	Operators   map[string]interface{} `json:"operators"`
}

type TemplateStrategy struct {
	Name       string                 `json:"name"`
	Cron       string                 `json:"cron"`
	Expression string                 `json:"expression"` //触发条件 表达式，检测变量名
	Script     string                 `json:"script"`     //javascript
	Operators  map[string]interface{} `json:"operators"`
}
