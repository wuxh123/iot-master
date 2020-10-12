package types

import "time"

type Model struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	H5          string    `json:"h5"`
	CreatedAt   time.Time `json:"created_at"`
}

type ModelBase struct {
	Id          int64     `json:"id"`
	ModelId     int64     `json:"model_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ModelTunnel struct {
	ModelBase       `xorm:"extends"`
	Protocol        string `json:"protocol"`
	ProtocolOpts    string `json:"protocol_opts"`
	PollingEnable   bool   `json:"polling_enable"`   //轮询
	PollingInterval int    `json:"polling_interval"` //轮询间隔 ms
	PollingCycle    int    `json:"polling_cycle"`    //轮询周期 s
}

type ModelVariable struct {
	ModelBase `xorm:"extends"`
	TunnelId  int    `json:"tunnel_id"`
	Type      string `json:"type"`
	Addr      string `json:"addr"`
	Alias     string `json:"alias"` //别名，用于编程
	Unit      string `json:"unit"`  //单位
	//应该不缩放，保留原始值？？？？
	Scale    float32 `json:"scale"` //倍率，比如一般是 整数÷10，得到
	Default  string  `json:"default"`
	Writable bool    `json:"writable"` //可写，用于输出（如开关）

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`
}

type ModelBatchResult struct {
	Id       int64  `json:"id"`
	BatchId  int64  `json:"batch_id"`
	Offset   int    `json:"offset"`
	Variable string `json:"variable"` //ModelVariable path
}

type ModelBatch struct {
	ModelBase `xorm:"extends"`
	TunnelId  int    `json:"tunnel_id"`
	Type      string `json:"type"`
	Addr      string `json:"addr"`
	Size      int    `json:"size"`

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`

	//结果解析 拆入子表
	//Results []ModelBatchResult `json:"results" xorm:"json"`
}

type ModelJob struct {
	ModelBase `xorm:"extends"`
	Cron      string `json:"cron"`
	Script    string `json:"script"` //javascript
}

type ModelStrategy struct {
	ModelBase `xorm:"extends"`
	Script    string `json:"script"` //javascript
}
