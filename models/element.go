package models

import "time"

type Address struct {
	Area      string `json:"area"`  //区域 类似 S I O Q WR ……
	Slave     uint8  `json:"slave"` //从站号 modbus
	Offset    uint16 `json:"offset"`
	ReadCode  uint8  `json:"read_code"`
	WriteCode uint8  `json:"write_code"`
}

type Element struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"` //来源

	Manufacturer string `json:"manufacturer"` //厂商
	Model        string `json:"model"`        //型号
	Version      string `json:"version"`      //版本

	ProtocolName string `json:"protocol_name"`
	ProtocolOpts string `json:"protocol_opts"`

	PollingEnable   bool `json:"polling_enable"`   //轮询
	PollingInterval int  `json:"polling_interval"` //轮询间隔 ms
	PollingCycle    int  `json:"polling_cycle"`    //轮询周期 s

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

type ElementVariable struct {
	Id        int64 `json:"id"`
	ElementId int64 `json:"element_id"`

	Address `xorm:"extends"` //地址

	Name  string `json:"name"`
	Alias string `json:"alias"` //别名，用于编程
	Type  string `json:"type"`
	Unit  string `json:"unit"` //单位

	Scale   float32 `json:"scale"`   //倍率，比如一般是 整数÷10，得到
	Correct float32 `json:"correct"` //校准

	Default  string `json:"default"`
	Writable bool   `json:"writable"` //可写，用于输出（如开关）

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}

type ElementBatch struct {
	Id        int64 `json:"id"`
	ElementId int64 `json:"element_id"`

	Address `xorm:"extends"`

	Size int `json:"size"`

	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
}
