package dtu

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
)

type Connect interface {
	onData(buf []byte)
	Send(buf []byte) (int, error)
	Close() error
}

type baseConnect struct {
	Error      string
	Serial     string
	RemoteAddr net.Addr

	channel *Channel
}

func (c *baseConnect) checkRegister(buf []byte) error {
	n := len(buf)
	if n < c.channel.Register.Length {
		return fmt.Errorf("register package is too short %d %s", n, string(buf[:n]))
	}
	serial := string(buf[:n])

	// 正则表达式判断合法性
	if c.channel.Register.Regex != "" {
		reg := regexp.MustCompile(`^` + c.channel.Register.Regex + `$`)
		match := reg.MatchString(serial)
		if !match {
			return fmt.Errorf("register package format error %s", serial)
		}
	}

	//按序号保存索引，供外部使用
	c.channel.serialIndexes.Store(serial, c)

	return nil
}

type connection struct {
	baseConnect
	conn net.Conn
}

func (c *connection) onData(buf []byte) {
	//检查注册包
	if c.channel.Register.Enable && c.Serial != "" {
		err := c.checkRegister(buf)
		if err != nil {
			_ = c.Close()
			return
		}

		//TODO 转发剩余内容

		return
	}

	//检查心跳包
	if c.channel.HeartBeat.Enable && bytes.Compare(c.channel.HeartBeat.Content, buf)==0{
		return
	}


	//TODO 内容转发

}

func (c *connection) Send(buf []byte) (int, error) {
	return c.conn.Write(buf)
}

func (c *connection) Close() error {
	return c.conn.Close()
}


func newConnection(conn net.Conn) *connection {
	return &connection{
		baseConnect: baseConnect{
			RemoteAddr: conn.RemoteAddr(),
		},
		conn: conn,
	}
}

type packetConnection struct {
	baseConnect
	conn net.PacketConn
}

func newPacketConnection(conn net.PacketConn, addr net.Addr) *packetConnection {
	return &packetConnection{
		baseConnect: baseConnect{
			RemoteAddr: addr,
		},
		conn: conn,
	}
}


func (c *packetConnection) onData(buf []byte) {
	//检查注册包
	if c.channel.Register.Enable && c.Serial != "" {
		err := c.checkRegister(buf)
		if err != nil {
			_ = c.Close()
			return
		}
		//TODO 转发剩余内容

		return
	}

	//检查心跳包
	if c.channel.HeartBeat.Enable && bytes.Compare(c.channel.HeartBeat.Content, buf)==0{
		return
	}

	//TODO 内容转发

}

func (c *packetConnection) Send(buf []byte) (int, error) {
	return c.conn.WriteTo(buf, c.RemoteAddr)
}

func (c *packetConnection) Close() error {
	return c.conn.Close()
}
