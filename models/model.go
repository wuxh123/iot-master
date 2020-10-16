package models

import "time"

type Model struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Disabled    bool      `json:"disabled"`
	Description string    `json:"description"`
	Origin      string    `json:"origin"` //模板ID
	Version     string    `json:"version"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
	Deployed    time.Time `json:"deployed"` //如果 deployed < updated，说明有更新，提示重新部署
}

type ModelBase struct {
	Id          int64     `json:"id"`
	ModelId     int64     `json:"model_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created" xorm:"created"`
	Updated     time.Time `json:"updated" xorm:"updated"`
}

type ModelAdapter struct {
	ModelBase `xorm:"extends"`

	ProtocolName string `json:"protocol_name"`
	ProtocolOpts string `json:"protocol_opts"`

	PollingEnable   bool `json:"polling_enable"`   //轮询
	PollingInterval int  `json:"polling_interval"` //轮询间隔 ms
	PollingCycle    int  `json:"polling_cycle"`    //轮询周期 s
}

type ModelVariable struct {
	ModelBase `xorm:"extends"`

	Address `xorm:"extends"`

	Type string `json:"type"`
	Unit string `json:"unit"` //单位

	Scale    float32 `json:"scale"` //倍率，比如一般是 整数÷10，得到
	Default  string  `json:"default"`
	Writable bool    `json:"writable"` //可写，用于输出（如开关）

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`
}

type ModelBatch struct {
	ModelBase `xorm:"extends"`

	Address `xorm:"extends"`

	Size int `json:"size"`

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`
}

type ModelJob struct {
	ModelBase `xorm:"extends"`

	Cron   string `json:"cron"`
	Script string `json:"script"` //javascript
}

type ModelStrategy struct {
	ModelBase `xorm:"extends"`

	Script string `json:"script"` //javascript
}
