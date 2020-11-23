package models

import (
	"time"
)

type Project struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid" storm:"unique"` //唯一码，自动生成
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Disabled bool `json:"disabled"`

	LinkId   int    `json:"link_id"`
	Protocol string `json:"protocol"`

	Variables  []ProjectVariable  `json:"variables"`
	Validators []ProjectValidator `json:"validators"`
	Functions  []ProjectFunction  `json:"functions"`
	Strategies []ProjectStrategy  `json:"strategies"`

	Created time.Time `json:"created" storm:"created"`
}

type ProjectVariable struct {
	Element  string `json:"element"` //uuid
	Variable `storm:"inline"`

	//TODO 添加采样周期
}

type ProjectValidator struct {
	Alert      string `json:"alert"`
	Expression string `json:"expression"` //表达式，检测变量名
}

type ProjectFunction struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"` //项目功能脚本唯一，供外部调用
	Description string                 `json:"description"`
	Script      string                 `json:"script"` //javascript
	Operators   map[string]interface{} `json:"operators"`
}

type ProjectStrategy struct {
	Name       string                 `json:"name"`
	Cron       string                 `json:"cron"`
	Expression string                 `json:"expression"` //触发条件 表达式，检测变量名
	Script     string                 `json:"script"`     //javascript
	Operators  map[string]interface{} `json:"operators"`
}
