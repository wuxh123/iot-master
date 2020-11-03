package base

//连接
type Link interface {
	Write(buf []byte) error
	Close() error

	Attach(peer Listener) error
	Detach() error

	Listen(listener Listener)
	//UnListen()
}

type Listener func([]byte)
