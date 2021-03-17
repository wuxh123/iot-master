package protocol

type Adapter interface {
	OnData(data []byte)

	Read(slave uint8, code uint8, offset uint16, size uint16) ([]byte, error)
	Write(slave uint8, code uint8, offset uint16, buf []byte) error
}
