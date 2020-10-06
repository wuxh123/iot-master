package types

import "time"

type ModelBase struct {
	Id          int       `json:"id" storm:"id,increment"`
	ModelId     int       `json:"model_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ModelTunnel struct {
	ModelBase       `storm:"inline"`
	LinkId          int    `json:"link_id"`
	Protocol        string `json:"protocol"`
	ProtocolOpts    string `json:"protocol_opts"`
	PollingEnable   bool   `json:"polling_enable"`   //轮询
	PollingInterval int    `json:"polling_interval"` //轮询间隔 ms
	PollingCycle    int    `json:"polling_cycle"`    //轮询周期 s
}

type ModelVariable struct {
	ModelBase `storm:"inline"`
	TunnelId  int    `json:"tunnel_id"`
	Type      string `json:"type"`
	Addr      string `json:"addr"`
	Alias     string `json:"alias"` //别名，用于编程
	Default   string `json:"default"`
	Writable  bool   `json:"writable"` //可写，用于输出（如开关）

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`
}

type ModelBatchResult struct {
	Offset   int    `json:"offset"`
	Variable string `json:"variable"` //ModelVariable path
}

type ModelBatch struct {
	ModelBase `storm:"inline"`
	TunnelId  int    `json:"tunnel_id"`
	Type      string `json:"type"`
	Addr      string `json:"addr"`
	Size      int    `json:"size"`

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`

	//结果解析
	Results []ModelBatchResult `json:"results"`
}

type ModelJob struct {
	ModelBase `storm:"inline"`
	Cron      string `json:"cron"`
	Script    string `json:"script"` //javascript
}

type ModelStrategy struct {
	ModelBase `storm:"inline"`
	Script    string `json:"script"` //javascript
}

type Model struct {
	Id          int       `json:"id" storm:"id,increment"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	H5          string    `json:"h5"`
	CreatedAt   time.Time `json:"created_at"`
}
