package model

import "errors"

//变量类型
type DataType uint8

const (
	DataTypeNone DataType = iota
	DataTypeBit
	DataTypeByte
	DataTypeWord
	DataTypeInteger
	DataTypeFloat
	DataTypeDouble
)

func ParseDataType(tp string) (DataType, error) {
	switch tp {
	case "bit":
		return DataTypeBit, nil
	case "byte":
		return DataTypeByte, nil
	case "word":
		return DataTypeWord, nil
	case "integer":
		return DataTypeInteger, nil
	case "float":
		return DataTypeFloat, nil
	case "double":
		return DataTypeDouble, nil
	default:
		return DataTypeNone, errors.New("未知类型")
	}
}

type Variable struct {
	//LinkId int
	Type   DataType
	Addr   string
	Value  interface{}
}
