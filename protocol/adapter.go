package protocol

import (
	"sync"
)

type AdapterListener interface {
	OnAdapterRead(addr string, buf []byte)
	OnAdapterWrite(addr string, size int)
	OnAdapterError(err error)
}

type Adapter interface {
	Listen(listener AdapterListener)

	Name() string
	Version() string

	Read(addr string, size int) error
	Write(addr string, buf []byte) error
}

//TODO 改为普通map
var adapters sync.Map

func Adapters() sync.Map {
	return adapters
}

//TODO 添加参数
func RegisterProtocol(name string, factory func() Adapter) {
	adapters.Store(name, factory)
}
