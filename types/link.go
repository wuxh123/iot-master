package types

//连接
type Link interface {
	Write(buf []byte) error
	Close() error

	Attach(peer OnDataFunc) error
	Detach() error

	Listen(listener LinkListener)
	//UnListen()
}

type LinkListener interface {
	OnData([]byte)
}

type OnDataFunc func([]byte)
