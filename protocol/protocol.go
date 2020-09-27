package protocol

import (
	"git.zgwit.com/zgwit/iot-admin/types"
)

type Protocol interface {
	Init()
	Read(typ types.DataType, addr string, size int) (cmd []byte, err error)
	Write(typ types.DataType, addr string, buf []byte) (cmd []byte, err error)
	Parse(buf []byte) (values []byte, err error)
}
