package base

import (
	"errors"
	"net"
)

type PackConn struct {
	net.PacketConn
	addr net.Addr
}

func (c *PackConn) Read(b []byte) (n int, err error) {
	return 0, errors.New("此接口不支持读")
}

func (c *PackConn) Write(b []byte) (n int, err error) {
	return c.WriteTo(b, c.addr)
}

func (c *PackConn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *PackConn) Close() error {
	//c.PacketConn = nil
	return nil
}
