package model

import "time"

type Element struct {
	Id     int64  `json:"id"`
	Origin string `json:"origin"` //平台UUID，自动生成

	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"` //厂商
	Model        string `json:"model"`        //型号
	Version      string `json:"version"`      //版本

	Variables []ElementVariable `json:"variables" xorm:"json"` //变量


	Created time.Time `json:"created" xorm:"created"`
}

//Modbus
// coil 线圈 离散输出 （1读多个、5写单个、15写多个）
// discrete 触点 离散输入 （2读多个）
// hold 保持寄存器（3读多个、6写单个、16写多个，--23读写多个--）
// input 输入寄存器（4读多个）

type ElementVariable struct {
	Variable
	Stretch uint16 `json:"stretch"` //扩展长度 默认0，如果大于1，自动在别名基础上添加数字后缀，比如 s s1 s2 s3 ...
}

type Variable struct {
	Name   string  `json:"name"`
	Alias  string  `json:"alias"`  //默认别名，用于编程
	Code   uint8   `json:"code"`   //功能码
	Offset uint16  `json:"offset"` //偏移
	Type   string  `json:"type"`
	Unit   string  `json:"unit"`  //单位
	Scale  float32 `json:"scale"` //倍率，比如一般是 整数÷10，得到

	//Default  string `json:"default"`
	ReadOnly bool `json:"read_only"` //只读
}
