package models

import "time"

type Element struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid" storm:"unique"` //唯一码，自动生成
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"` //来源

	Manufacturer string `json:"manufacturer"` //厂商
	Model        string `json:"model"`        //型号
	Version      string `json:"version"`      //版本

	//Variables []ElementVariable `json:"variables"`
	Created time.Time `json:"created" storm:"created"`
}

//Modbus Area
// discrete 离散输入 触点（2读多个）
// coil 离散输出 线圈（1读多个、5写单个、15写多个）
// input 输入寄存器（4读多个）
// hold 保持寄存器（3读多个、6写单个、16写多个，--23读写多个--）

type ElementVariable struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Alias  string `json:"alias"` //别名，用于编程
	Area   string `json:"area"`
	Offset uint16 `json:"offset"`
	Type   string `json:"type"`
	Unit   string `json:"unit"` //单位

	Scale float32 `json:"scale"` //倍率，比如一般是 整数÷10，得到

	Default  string `json:"default"`
	ReadOnly bool   `json:"read_only"` //只读

	Created time.Time `json:"created" storm:"created"`
}
