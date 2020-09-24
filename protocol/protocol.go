package protocol

import "git.zgwit.com/zgwit/iot-admin/model"

type Protocol interface {
	Init()
	Read(typ model.DataType, addr string, length int) ([]byte, error)
	Write(typ model.DataType, addr string, buf []byte) ([]byte, error)
	Parse(buf []byte) ([]byte, error)
}
