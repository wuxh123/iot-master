package cool

import "git.zgwit.com/zgwit/iot-admin/models"

//通道
type Tunnel interface {
	Close() error
	GetTunnel() *models.ModelTunnel
}

//连接
type Link interface {
	Write(buf []byte) error
	Close() error

	Attach(link Link) error
	Detach() error
}
