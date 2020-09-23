package protocol



type Protocol interface {
	Init()
	Read(addr string, length int) ([]byte, error)
	Write(addr string, buf []byte) ([]byte, error)
	Parse(buf []byte) ([]byte, error)
}
