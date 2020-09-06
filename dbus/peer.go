package dbus

import (
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/packet"
	"log"
)


type Peer struct {
	baseClient

	link *dtu.Link
}

func (p *Peer) Handle(msg *packet.Packet) {
	switch msg.Type {
	case packet.TypeConnect:
		_ = p.SendError("duplicate register")
		_ = p.CLose()
	case packet.TypeHeartBeak:
	case packet.TypePing:
		_ = p.Send(&packet.Packet{Type: packet.TypePong})
	case packet.TypeTransfer:
		p.handleTransfer(msg)
	default:
		log.Println("unknown command", msg)
	}
}

func (p *Peer) handleTransfer(msg *packet.Packet) {
	//_, _ = p.link.Send(msg.Data)
	//TODO 判断link是否为空
	//TODO 如果link断线，缓存 数据包

}
