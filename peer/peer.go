package peer

import (
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/packet"
	"log"
	"net"
	"time"
)

type Peer struct {
	Key  string
	conn net.Conn

	link *dtu.Link
}

func (p *Peer) CLose() {
	p.conn = nil
	//TODO 从列表中删除，从link中删除
}

func (p *Peer) Send(msg *packet.Packet) error {
	return p.Write(msg.Encode())
}

func (p *Peer) Write(b []byte) error {
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

func (p *Peer) handle(msg *packet.Packet) {
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

func (p *Peer) handleConnect(msg *packet.Packet) {
	p.Key = string(msg.Data)
	//TODO 检查Key，获取对应的 通道，连接号
	//TODO 在link中设置peer

}

func (p *Peer) handleTransfer(msg *packet.Packet) {
	_, _ = p.link.Send(msg.Data)
	//TODO 判断link是否为空
	//TODO 如果link断线，缓存 数据包

}
