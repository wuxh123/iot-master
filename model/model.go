package model

//变量类型
type DataType uint8

const (
	DataTypeBit DataType = iota
	DataTypeByte
	DataTypeWord
	DataTypeInteger
	DataTypeFloat
	DataTypeDouble
)
