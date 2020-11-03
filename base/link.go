package base

//连接
type Link interface {
	Write(buf []byte) error
	Close() error

	Attach(peer OnDataFunc) error
	Detach() error

	Listen(listener OnDataFunc)
	//UnListen()
}

type OnDataFunc func([]byte)
