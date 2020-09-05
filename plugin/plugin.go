package plugin

import (
	"github.com/zgwit/dtu-admin/packet"
	"log"
	"net"
	"time"
)

type Plugin struct {
	conn net.Conn
}

func (p *Plugin) CLose() {
	p.conn = nil
}

func (p *Plugin) Send(msg *packet.Packet) error {
	return p.Write(msg.Encode())
}

func (p *Plugin) Write(b []byte) error {
	n, e := p.conn.Write(b)
	if e != nil {
		return e
	}

	//继续发送????，没太大必要？
	if n < len(b) {
		log.Println("发送不完整")
		time.Sleep(time.Millisecond * 100)
		_, _ = p.conn.Write(b[n:])
	}

	return nil
}

func (p *Plugin) handle(msg *packet.Packet) {
	switch msg.Type {
	case packet.TypeConnect:
		p.handleConnect(msg)
	case packet.TypeHeartBeak:
	case packet.TypePing:
		_ = p.Send(&packet.Packet{Type: packet.TypePong})
	case packet.TypeTransfer:
		p.handleTransfer(msg)
	default:
		log.Println("unknown command", msg)
	}
}

func (p *Plugin) handleConnect(msg *packet.Packet) {
	//TODO 根据appkey, secret校验身份，注册插件到对应通道和链接上

}

func (p *Plugin) handleTransfer(msg *packet.Packet) {
	//TODO 找到对应链接，发送之
	//TODO 使用 int64 还是 int32
	//p.link.Send(msg.Data[8:])

}
