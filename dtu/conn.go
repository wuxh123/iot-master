package dtu

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"regexp"
)

type Connection struct {
	Error      string
	Serial     string
	RemoteAddr net.Addr

	conn interface{}

	channel *Channel
}

func (c *Connection) checkRegister(buf []byte) error {
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
	connections.Store(serial, c)

	return nil
}

func (c *Connection) onData(buf []byte) {
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
	if c.channel.HeartBeat.Enable && bytes.Compare(c.channel.HeartBeat.Content, buf) == 0 {
		return
	}

	//TODO 内容转发，暂时直接回复
	_, _ = c.Send(buf)

}

func (c *Connection) Send(buf []byte) (int, error) {

	if conn, ok := c.conn.(net.Conn); ok {
		return conn.Write(buf)
	}
	if conn, ok := c.conn.(net.PacketConn); ok {
		return conn.WriteTo(buf, c.RemoteAddr)
	}
	return 0, errors.New("错误的链接类型")
}


func (c *Connection) Close() error {
	return c.conn.(net.Conn).Close()
}

func newConnection(conn net.Conn) *Connection {
	return &Connection{
		RemoteAddr: conn.RemoteAddr(),
		conn:       conn,
	}
}


func newPacketConnection(conn net.PacketConn, addr net.Addr) *Connection {
	return &Connection{
			RemoteAddr: addr,
		conn: conn,
	}
}
