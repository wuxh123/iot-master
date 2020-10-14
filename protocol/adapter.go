package protocol

import (
	"git.zgwit.com/zgwit/iot-admin/models"
	"sync"
)

type AdapterListener interface {
	OnAdapterRead(addr *models.Address, buf []byte)
	OnAdapterWrite(addr *models.Address, size int)
	OnAdapterError(err error)
}

type Adapter interface {
	Listen(listener AdapterListener)

	Name() string
	Version() string

	Read(addr *models.Address, size int) error
	Write(addr *models.Address, buf []byte) error
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
