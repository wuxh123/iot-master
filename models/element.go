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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"` //来源
	UUID        string `json:"uuid"`   //唯一码，自动生成

	Manufacturer string `json:"manufacturer"` //厂商
	Model        string `json:"model"`        //型号
	Version      string `json:"version"`      //版本

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}

//Modbus Area
// discrete 离散输入 触点（2读多个）
// coil 离散输出 线圈（1读多个、5写单个、15写多个）
// input 输入寄存器（4读多个）
// hold 保持寄存器（3读多个、6写单个、16写多个，--23读写多个--）

type ElementVariable struct {
	ID        int `json:"id"`
	ElementId int `json:"element_id"`

	Name string `json:"name"`
	//Alias  string `json:"alias"` //别名，用于编程
	Area   string `json:"area"`
	Offset uint16 `json:"offset"`
	Type   string `json:"type"`
	Unit   string `json:"unit"` //单位

	Scale   float32 `json:"scale"`   //倍率，比如一般是 整数÷10，得到

	Default  string `json:"default"`
	ReadOnly bool   `json:"read_only"` //只读

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}
