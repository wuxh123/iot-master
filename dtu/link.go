package dtu

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/asdine/storm/v3/q"
	"github.com/zgwit/dtu-admin/storage"
	"github.com/zgwit/dtu-admin/model"
	"log"
	"net"
	"regexp"
	"time"
)

type Link struct {
	ID int

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

	//配置序列号
	l.Serial = serial

	//查找数据库同通道，同序列号链接，更新数据库中 serial online
	db := storage.DB("link")
	var link model.Link
	err := db.Select(q.Eq("Channel", l.channel.ID), q.Eq("Serial", serial)).First(&link)
	if err == nil {
		//检查工作状态，如果同序号连接还在正常通讯，则关闭当前连接，回复：Duplicate register

		//更新客户端地址，
		err = storage.DB("link").Update(&model.Link{
			ID:     link.ID,
			Addr:   l.RemoteAddr.String(),
			Online: time.Now(),
		})
		if err != nil {
			log.Println(err)
		}
	} else {
		link = model.Link{
			Serial:  serial,
			Addr:    l.RemoteAddr.String(),
			Channel: l.channel.ID,
			Online:  time.Now(),
			Created: time.Now(),
		}
		err = storage.DB("link").Save(&link)
		if err != nil {
			log.Println(err)
		}
		l.ID = link.ID
	}

	return nil
}

func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

	//检查注册包
	if l.channel.Register.Enable && l.Serial == "" {
		err := l.checkRegister(buf)
		if err != nil {
			log.Println(err)
			_ = l.Close()
			return
		}

		//TODO 转发剩余内容

		return
	}

	hb := l.channel.HeartBeat
	//检查心跳包, 判断上次收发时间，是否已经过去心跳间隔
	if hb.Enable && time.Now().Sub(l.lastTime) > time.Second*time.Duration(hb.Interval) {
		var b []byte
		if hb.IsHex {
			var e error
			b, e = hex.DecodeString(hb.Content)
			if e != nil {
				log.Println(e)
			}
		} else {
			b = []byte(hb.Content)
		}
		if bytes.Compare(b, buf) == 0 {
			return
		}
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
