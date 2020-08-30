package plugin

import (
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	conn net.Conn
}

func (c *Client) Send(msg *Message) error {
	return c.Write(msg.Encode())
}

func (c *Client) Write(b []byte) error {
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

func (c *Client) handle(msg *Message) {
	log.Println("handle message", msg)

	ss := strings.Split(msg.cmd, "/")

	cmd := ss[0]
	args := ss[1:]

	switch cmd {
	case "register":
		c.handleRegister(args)
	case "subscribe":
		c.handleSubscribe(args)
	case "unsubscribe":
		c.handleUnsubscribe(args)
	case "transfer":
		c.handleTransfer(args)
	//case "event": // 应该只有插件处理
	default:
		log.Println("unknown command", msg)
	}
}

func (c *Client) handleRegister(args []string) {

}

// subscribe/:channel/:id
func (c *Client) handleSubscribe(args []string) {

}

// unsubscribe/:channel/:id
func (c *Client) handleUnsubscribe(args []string) {

}

//开始 transfer/:channel/:id/start（非必需调用）
//透传 transfer/:channel/:id,
//停止 transfer/:channel/:id/stop
func (c *Client) handleTransfer(args []string) {

}

