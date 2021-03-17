package types

import "iot-master/model"

//通道
type Tunnel interface {
	Open() error
	Close() error
	GetModel() *model.Tunnel
	GetLink(id int64) (Link, error)
}
