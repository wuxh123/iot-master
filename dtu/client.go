package dtu

import (
	"errors"
	"net"
)

type Client struct {
	Net  string
	Addr string
	Err  string

	Serial string

	conn net.Conn

	HeartBeatConf HeartBeatConf
}

func NewClient(net string, addr string) *Client {
	return &Client{
		Net:  net,
		Addr: addr,
	}
}

func (c *Client) Open() error {
	var err error
	c.conn, err = net.Dial(c.Net, c.Addr)
	if err != nil {
		return err
	}
	go c.receive()
	return nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return errors.New("client closed")
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *Client) receive() {

}
