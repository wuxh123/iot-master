package dbus

import (
	"github.com/zgwit/dtu-admin/packet"
	"log"
	"net"
	"time"
)

type baseClient struct {
	conn net.Conn
}

func (p *baseClient) CLose() error {
	if p.conn == nil {
		return nil
	}
	err := p.conn.Close()
	p.conn = nil
	return err
}

func (p *baseClient) Send(msg *packet.Packet) error {
	return p.Write(msg.Encode())
}

func (p *baseClient) Write(b []byte) error {
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

func (p *baseClient) SendError(err string) error {
	return p.Send(&packet.Packet{
		Type:   packet.TypeError,
		Status: 0,
		Data:   []byte(err),
	})
}

func (p *baseClient) Disconnect(info string) error {
	return p.Send(&packet.Packet{
		Type:   packet.TypeDisconnect,
		Status: 0,
		Data:   []byte(info),
	})
}