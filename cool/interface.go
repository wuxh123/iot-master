package cool

//通道
type Tunnel interface {
	Close() error
}

//连接
type Link interface {
	Write(buf []byte) error
	Close() error

	Attach(link Link) error
	Detach() error
}
