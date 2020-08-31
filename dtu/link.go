package dtu

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

type Link struct {
	ID int64

	Error      string
	Serial     string
	RemoteAddr net.Addr

	Rx int
	Tx int

	conn interface{}

	lastTime time.Time

	channel *Channel

}

func (l *Link) checkRegister(buf []byte) error {
	n := len(buf)
	if n < l.channel.Register.Length {
		return fmt.Errorf("register package is too short %d %s", n, string(buf[:n]))
	}
	serial := string(buf[:n])

	// 正则表达式判断合法性
	if l.channel.Register.Regex != "" {
		reg := regexp.MustCompile(`^` + l.channel.Register.Regex + `$`)
		match := reg.MatchString(serial)
		if !match {
			return fmt.Errorf("register package format error %s", serial)
		}
	}

	//按序号保存索引，供外部使用
	connections.Store(serial, l)

	return nil
}

func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

	//检查注册包
	if l.channel.Register.Enable && l.Serial != "" {
		err := l.checkRegister(buf)
		if err != nil {
			log.Println(err)
			_ = l.Close()
			return
		}

		//TODO 转发剩余内容

		return
	}

	//检查心跳包
	if l.channel.HeartBeat.Enable && bytes.Compare(l.channel.HeartBeat.Content, buf) == 0 {
		//TODO 判断上次收发时间，是否已经过去心跳间隔
		return
	}

	//TODO 内容转发，暂时直接回复
	_, _ = l.Send(buf)

}

func (l *Link) Send(buf []byte) (int, error) {
	l.Tx += len(buf)
	l.lastTime = time.Now()

	if conn, ok := l.conn.(net.Conn); ok {
		return conn.Write(buf)
	}
	if conn, ok := l.conn.(net.PacketConn); ok {
		return conn.WriteTo(buf, l.RemoteAddr)
	}
	return 0, errors.New("错误的链接类型")
}

func (l *Link) Close() error {
	return l.conn.(net.Conn).Close()
}

func newConnection(conn net.Conn) *Link {
	return &Link{
		RemoteAddr: conn.RemoteAddr(),
		conn:       conn,
	}
}

func newPacketConnection(conn net.PacketConn, addr net.Addr) *Link {
	return &Link{
		RemoteAddr: addr,
		conn:       conn,
	}
}
