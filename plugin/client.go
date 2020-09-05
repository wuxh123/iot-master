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

func (c *Plugin) Send(msg *packet.Packet) error {
	return c.Write(msg.Encode())
}

func (c *Plugin) Write(b []byte) error {
	n, e := c.conn.Write(b)
	if e != nil {
		return e
	}

	//继续发送????，没太大必要？
	if n < len(b) {
		log.Println("发送不完整")
		time.Sleep(time.Millisecond * 100)
		_, _ = c.conn.Write(b[n:])
	}

	return nil
}

func (c *Plugin) handle(msg *packet.Packet) {
	//log.Println("handle message", msg)

	switch msg.Type {
	case packet.TypeConnect:
		c.handleRegister(msg)
	case packet.TypeHeartBeak:
	case packet.TypePing:
		_ = c.Send(&packet.Packet{Type: packet.TypePong})
	case packet.TypeTransfer:
		c.handleTransfer(msg)
	default:
		log.Println("unknown command", msg)
	}
}

func (c *Plugin) handleRegister(msg *packet.Packet) {
	//TODO 根据appkey, secret校验身份，注册插件到对应通道和链接上

}

func (c *Plugin) handleTransfer(msg *packet.Packet) {
	//TODO 找到对应链接，发送之

}
