package model

import (
	"git.zgwit.com/zgwit/iot-admin/model/json"
	"github.com/robertkrimen/otto"
	"sync"
)
import "errors"

//变量类型
type DataType uint8

const (
	DataTypeNone DataType = iota
	DataTypeBit
	DataTypeByte
	DataTypeUint8
	DataTypeUint16
	DataTypeUint32
	DataTypeUint64
	DataTypeInt8
	DataTypeInt16
	DataTypeInt32
	DataTypeInt64
	DataTypeFloat32
	DataTypeFloat64
	DataTypeFloat128
)

func (dt DataType) ToString() string {
	switch dt {
	case DataTypeBit:
		return "bit"
	default:
		return "unknown"
	}
}

func ParseDataType(tp string) (DataType, error) {
	switch tp {
	case "bit":
		return DataTypeBit, nil
	case "byte":
		return DataTypeByte, nil
		//TODO 填充类型
	default:
		return DataTypeNone, errors.New("未知类型")
	}
}

type Variable struct {
	//LinkId int
	Type  DataType
	Addr  string
	Value interface{}
}

type Model struct {
	vm *otto.Otto

	variables sync.Map

	jobs       []json.Job
	strategies []json.Strategy
}

func NewModel() *Model {
	return &Model{
		vm:         otto.New(),
		jobs:       make([]json.Job, 0),
		strategies: make([]json.Strategy, 0),
	}
}
