package base

//连接
type Link interface {
	Listen(listener LinkListener)

	Write(buf []byte) error
	Close() error

	Attach(link Link) error
	Detach() error
}

type LinkListener interface {
	OnLinkData(buf []byte)
	//OnLinkError(err error)
	//OnLinkerClose()
}
