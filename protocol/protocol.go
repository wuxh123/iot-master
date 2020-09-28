package protocol

import (
	"git.zgwit.com/zgwit/iot-admin/interfaces"
	"git.zgwit.com/zgwit/iot-admin/types"
	"sync"
)

type AdapterListener interface {
	OnAdapterRead(addr string, typ types.DataType, buf []byte)
	OnAdapterWrite(addr string, typ types.DataType, size int)
	OnAdapterError(err error)
}

type Adapter interface {
	Listen(listener AdapterListener)

	Name() string
	Version() string

	Read(addr string, typ types.DataType, size int) (err error)
	Write(addr string, typ types.DataType, buf []byte) (err error)
}

//TODO 改为普通map
var adapters sync.Map

func Adapters() sync.Map {
	return adapters
}

//TODO 添加参数
func RegisterProtocol(name string, factory func(linker interfaces.Linker) Adapter) {
	adapters.Store(name, factory)
}
