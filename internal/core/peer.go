package core

import "net"

type Peer struct {
	conn net.Conn
	link *Link
}

func NewPeer(conn net.Conn, link *Link) *Peer {
	peer := &Peer{
		conn: conn,
		link: link,
	}
	link.peer = peer
	return peer
}

func (p *Peer) receive() {
	buf := make([]byte, 1024)
	for {
		n, e := p.conn.Read(buf)
		if e != nil {
			break
		}
		_ = p.link.Write(buf[:n])
		//TODO 使用协议
	}
	_ = p.Close()
}

func (p *Peer) Send(buf []byte) error {
	return nil
}

func (p *Peer) Close() error {
	p.link.peer = nil
	return p.conn.Close()
}
