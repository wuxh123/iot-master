package base

import "git.zgwit.com/iot/dtu-admin/packet"

type Tunnel interface {
	CLose() error
	Send(msg *packet.Packet) error
	Write(b []byte) error
	Handle(msg *packet.Packet)
}


