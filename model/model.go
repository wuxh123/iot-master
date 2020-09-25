package model

import (
	"git.zgwit.com/zgwit/iot-admin/model/json"
	"github.com/robertkrimen/otto"
	"strings"
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
const _typeNames = ",bit,byte,uint8,uin16,uin32,uint64,int8,int16,int32,int64,float32,float64,float128"
var typeNames []string

func init() {
	typeNames = strings.Split(_typeNames, ",")
}

func (dt DataType) ToString() string {
	return typeNames[dt]
}

func ParseDataType(tp string) (DataType, error) {
	for i, t := range typeNames {
		if t == tp {
			return DataType(i), nil
		}
	}
	return DataTypeNone, errors.New("未知类型")
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
