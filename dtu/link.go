package dtu

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/model"
	"log"
	"net"
	"regexp"
	"time"
)

type Link struct {
	Id int64

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
	if n < l.channel.RegisterMin {
		return fmt.Errorf("register package is too short %d %s", n, string(buf[:n]))
	}
	serial := string(buf[:n])
	if n > l.channel.RegisterMax {
		serial = string(buf[:l.channel.RegisterMax])
	}

	// 正则表达式判断合法性
	if l.channel.RegisterRegex != "" {
		reg := regexp.MustCompile(`^` + l.channel.RegisterRegex + `$`)
		match := reg.MatchString(serial)
		if !match {
			return fmt.Errorf("register package format error %s", serial)
		}
	}

	//配置序列号
	l.Serial = serial

	//查找数据库同通道，同序列号链接，更新数据库中 addr online
	var link model.Link
	has, err := db.Engine.Where("channel_id=?", l.channel.Id).And("serial=?", serial).Get(&link)
	if err != nil {
		return err
	}
	if has {
		//TODO 检查工作状态，如果同序号连接还在正常通讯，则关闭当前连接，回复：Duplicate register

		//更新客户端地址，
		link.Addr = l.RemoteAddr.String()
		link.Online = time.Now()
		_, err := db.Engine.ID(link.Id).Cols("addr", "online").Update(link)
		if err != nil {
			return err
		}
	} else {
		link = model.Link{
			Serial:    serial,
			Addr:      l.RemoteAddr.String(),
			ChannelId: l.channel.Id,
			Online:    time.Now(),
			Created:   time.Now(),
		}
		_, err := db.Engine.Insert(&link)
		if err != nil {
			return err
		}
		l.Id = link.Id
	}

	//处理剩余内容
	if n > l.channel.RegisterMax {
		l.onData(buf[l.channel.RegisterMax:])
	}

	return nil
}

func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

	//检查注册包
	if l.channel.RegisterEnable && l.Serial == "" {
		err := l.checkRegister(buf)
		if err != nil {
			log.Println(err)
			_, _ = l.Send([]byte(err.Error()))
			_ = l.Close()
			return
		}
		return
	}

	//检查心跳包, 判断上次收发时间，是否已经过去心跳间隔
	if l.channel.HeartBeatEnable && time.Now().Sub(l.lastTime) > time.Second*time.Duration(l.channel.HeartBeatInterval) {
		var b []byte
		if l.channel.HeartBeatIsHex {
			var e error
			b, e = hex.DecodeString(l.channel.HeartBeatContent)
			if e != nil {
				log.Println(e)
			}
		} else {
			b = []byte(l.channel.HeartBeatContent)
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
