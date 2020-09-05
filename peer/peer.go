package peer

import (
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/packet"
	"log"
	"net"
)

type Peer struct {
	Key  string
	conn net.Conn

	//发送缓存
	cache [][]byte

	link   *dtu.Link
	parser packet.Parser
}

func NewPeer(key string, conn net.Conn) *Peer {
	return &Peer{
		Key:  key,
		conn: conn,
	}
}

func (p *Peer) CLose() {
	p.conn = nil
	//TODO 从列表中删除，从link中删除
}

func (p *Peer) Send(msg *packet.Packet) error {
	_, err := p.conn.Write(msg.Encode())
	return err
}

func (p *Peer) receive() {
	buf := make([]byte, 1024)
	for p.conn != nil {
		n, e := p.conn.Read(buf)
		if e != nil {
			log.Println(e)
			break
		}
		packs := p.parser.Parse(buf[:n])
		for _, pack := range packs {
			p.handle(pack)
		}
	}
}

func (p *Peer) handle(msg *packet.Packet) {
	switch msg.Type {
	//
	case packet.TypeConnect:
		p.Key = string(msg.Data)
	//TODO 检查Key，获取对应的 通道，连接号
	//TODO 在link中设置peer

	case packet.TypeHeartBeak:
	case packet.TypePing:
		_ = p.Send(&packet.Packet{
			Type: packet.TypePong,
			Data: nil,
		})
	case packet.TypeTransfer:
		_, _ = p.link.Send(msg.Data)
		//TODO 判断link是否为空
		//TODO 如果link断线，缓存 数据包
	}
}
