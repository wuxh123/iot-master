package protocol

import (
	"iot-master/base"
)

type Adapter interface {
	Attach(link base.Link)

	Read(slave uint8, code uint8, offset uint16, size uint16) ([]byte, error)
	Write(slave uint8, code uint8, offset uint16, buf []byte) error
}


