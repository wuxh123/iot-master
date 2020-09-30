package core

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type Peer struct {
	conn *websocket.Conn
	link *Link
}

type PeerPacket struct {
	cmd     string
	payload []byte
}

func NewPeer(conn *websocket.Conn, link *Link) *Peer {
	peer := &Peer{
		conn: conn,
		link: link,
	}
	link.peer = peer
	return peer
}

func (p *Peer) Receive() {
	for {
		var pack PeerPacket
		if err := websocket.JSON.Receive(p.conn, &pack); err != nil {
			fmt.Println(err)
			break
		}

		switch pack.cmd {
		case "transfer":
			_ = p.link.Write(pack.payload)
		}
	}
	_ = p.Close()
}

func (p *Peer) Send(buf []byte) error {
	pack := &PeerPacket{
		cmd:     "transfer",
		payload: buf,
	}
	return websocket.Message.Send(p.conn, pack)
}

func (p *Peer) Close() error {
	p.link.peer = nil
	return p.conn.Close()
}
