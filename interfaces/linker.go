package interfaces

type LinkerListener interface {
	OnLinkerData(buf []byte)
	OnLinkerError(err error)
	OnLinkerClose()
}

type Linker interface {
	Listen(listener LinkerListener)
	Write(buf []byte) error
	Close() error
}

