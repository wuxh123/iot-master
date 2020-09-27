package protocol

import (
	"git.zgwit.com/zgwit/iot-admin/types"
	"net"
	"sync"
)

type Listener interface {
	onRead(addr string, typ types.DataType, buf []byte)
	onWrite(addr string, typ types.DataType, size int)
	onError(err error)
}

type Protocol interface {
	//TODO 其他参数
	Init(listener Listener)
	Name() string

	Read(addr string, typ types.DataType, size int) (err error)
	Write(addr string, typ types.DataType, buf []byte) (err error)
}

var protocols sync.Map

func RegisterProtocol(name string, factory func(conn net.Conn) Protocol)  {

}
