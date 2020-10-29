package protocol

import (
	"sync"
)

type Address struct {
	Slave  uint8
	Offset uint16
	Area   string
}

type AdapterListener interface {
	OnAdapterRead(addr *Address, buf []byte)
	OnAdapterWrite(addr *Address, size int)
	OnAdapterError(err error)
}

type Adapter interface {
	Listen(listener AdapterListener)

	Name() string
	Version() string

	Read(addr *Address, size int) error
	Write(addr *Address, buf []byte) error
}

//可以改为普通map
var adapters sync.Map

func Adapters() {
	//return adapters
}

//TODO 添加参数
func RegisterProtocol(name string, factory func() Adapter) {
	adapters.Store(name, factory)
}
