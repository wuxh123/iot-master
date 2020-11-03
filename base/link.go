package base

//连接
type Link interface {
	Request(buf []byte) ([]byte, error)

	Write(buf []byte) error
	Close() error

	Attach(link Link) error
	Detach() error
}
