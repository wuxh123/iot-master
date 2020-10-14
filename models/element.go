package models

import "time"

type Element struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Origin      string    `json:"origin"` //模板ID
	Version     string    `json:"version"`
	Protocol    string    `json:"protocol"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

type ElementVariable struct {
	Id          int64     `json:"id"`
	ElementId   int64     `json:"element_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`

	Type  string `json:"type"`
	Addr  string `json:"addr"`
	Alias string `json:"alias"` //别名，用于编程
	Unit  string `json:"unit"`  //单位
	//应该不缩放，保留原始值？？？？
	Scale    float32 `json:"scale"` //倍率，比如一般是 整数÷10，得到
	Default  string  `json:"default"`
	Writable bool    `json:"writable"` //可写，用于输出（如开关）
}
