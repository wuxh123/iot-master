package core

import "net"

type Peer struct {
	conn net.Conn
	link *Link
}

func PeerAccept(conn net.Conn) {

}

func (p *Peer) Send(buf []byte) error {
	return nil
}

func (p *Peer) Close(buf []byte) error {
	return nil
}
