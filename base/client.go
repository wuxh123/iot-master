package base

import "github.com/zgwit/dtu-admin/packet"

type Client interface {
	CLose() error
	Send(msg *packet.Packet) error
	Write(b []byte) error
	Handle(msg *packet.Packet)
}


