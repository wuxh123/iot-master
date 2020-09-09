package dbus

import (
	"fmt"
	"github.com/zgwit/dtu-admin/packet"
)

type Plugin struct {
	baseClient
}

func (p *Plugin) Handle(msg *packet.Packet) {
	switch msg.Type {
	case packet.TypeConnect:
		_ = p.Disconnect("duplicate register")
		_ = p.CLose()
	case packet.TypeHeartBeak:
	case packet.TypePing:
		_ = p.Send(&packet.Packet{Type: packet.TypePong})
	case packet.TypeTransfer:
		p.handleTransfer(msg)
	default:
		_ = p.SendError(fmt.Sprintf("unknown command %d", msg.Type))
	}
}

func (p *Plugin) handleTransfer(msg *packet.Packet) {
	//TODO 找到对应链接，发送之
	//TODO 使用 int 还是 int32
	//p.link.Send(msg.Data[8:])

}
