package dtu

import (
	"github.com/zgwit/dtu-admin/packet"
	"log"
	"net"
)

type Peer struct {
	Key  string
	conn net.Conn

	parser packet.Parser
}

func NewPeer(key string, conn net.Conn) *Peer {
	return &Peer{
		Key:  key,
		conn: conn,
	}
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
	}
}
