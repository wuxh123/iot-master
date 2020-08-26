package dtu

import (
	"errors"
	"net"
)

type Client struct {
	Net  string
	Addr string

	conn net.Conn
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

func (s *Client) Close() error {
	if s.conn == nil {
		return errors.New("client closed")
	}
	err := s.conn.Close()
	s.conn = nil
	return err
}

func (s *Client) receive() {

}