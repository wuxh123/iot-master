package base

//连接
type Link interface {
	Write(buf []byte) error
	Close() error

	Attach(link Link) error
	Detach() error
}

