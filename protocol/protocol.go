package protocol

import (
	"git.zgwit.com/zgwit/iot-admin/types"
	"sync"
)

type LinkerListener interface {
	onLinkerData(buf []byte)
	onLinkerError(err error)
	onLinkerClose()
}

type Linker interface {
	Listen(listener LinkerListener) error
	Write(buf []byte)
}

type AdapterListener interface {
	onAdapterRead(addr string, typ types.DataType, buf []byte)
	onAdapterWrite(addr string, typ types.DataType, size int)
	onAdapterError(err error)
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
func RegisterProtocol(name string, factory func(linker Linker) Adapter) {
	adapters.Store(name, factory)
}
