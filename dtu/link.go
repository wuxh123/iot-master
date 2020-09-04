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
	model.Link

	registerChecked bool

	//RemoteAddr net.Addr

	Rx int
	Tx int

	conn net.Conn

	lastTime time.Time

	channel Channel
}

func (l *Link) checkRegister(buf []byte) error {
	ch := l.channel.GetChannel()

	n := len(buf)
	if n < ch.RegisterMin {
		return fmt.Errorf("register package is too short %d %s", n, string(buf[:n]))
	}
	serial := string(buf[:n])
	if ch.RegisterMax > 0 && ch.RegisterMax >= ch.RegisterMin && n > ch.RegisterMax {
		serial = string(buf[:ch.RegisterMax])
	}

	// 正则表达式判断合法性
	if ch.RegisterRegex != "" {
		reg := regexp.MustCompile(`^` + ch.RegisterRegex + `$`)
		match := reg.MatchString(serial)
		if !match {
			return fmt.Errorf("register package format error %s", serial)
		}
	}

	//配置序列号
	l.Serial = serial

	//查找数据库同通道，同序列号链接，更新数据库中 addr online
	var link model.Link
	has, err := db.Engine.Where("channel_id=?", ch.Id).And("serial=?", serial).Get(&link)
	if err != nil {
		return err
	}
	if has {
		lnk, _ := l.channel.GetLink(link.Id)
		if lnk != nil {
			//如果同序号连接还在正常通讯，则关闭当前连接
			if lnk.conn != nil {
				return fmt.Errorf("duplicate serial %s", serial)
			}

			//复制有用的历史数据
			l.Rx = lnk.Rx
			l.Tx = lnk.Tx

			//复制watcher
		}

		l.Id = link.Id
		l.Name = link.Name
		l.Serial = link.Serial
	}

	//保存链接
	l.channel.StoreLink(l)

	//处理剩余内容
	if ch.RegisterMax > 0 && n > ch.RegisterMax {
		l.onData(buf[ch.RegisterMax:])
	}

	return nil
}

func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

	ch := l.channel.GetChannel()

	//检查注册包（只有服务端是检测）
	if !l.registerChecked && ch.RegisterEnable && ch.Role == "server" {
		err := l.checkRegister(buf)
		if err != nil {
			log.Println(err)
			_, _ = l.Send([]byte(err.Error()))
			_ = l.Close()
			return
		}
		l.registerChecked = true
		return
	}

	//检查心跳包, 判断上次收发时间，是否已经过去心跳间隔
	if ch.HeartBeatEnable && time.Now().Sub(l.lastTime) > time.Second*time.Duration(ch.HeartBeatInterval) {
		var b []byte
		if ch.HeartBeatIsHex {
			var e error
			b, e = hex.DecodeString(ch.HeartBeatContent)
			if e != nil {
				log.Println(e)
			}
		} else {
			b = []byte(ch.HeartBeatContent)
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

	return l.conn.Write(buf)
}

func (l *Link) Close() error {
	if l.conn == nil {
		return errors.New("连接已经关闭")
	}
	err := l.conn.Close()
	l.conn = nil
	if err != nil {
		return err
	}
	l.Online = false
	_, err = db.Engine.ID(l.Id).Cols("online").Update(&l.Link)
	return err
}

func (l *Link) storeError(err error) error {
	l.Error = err.Error()
	_, err = db.Engine.ID(l.Id).Cols("error").Update(&l.Link)
	return err
}

func newLink(ch Channel, conn net.Conn) *Link {
	c := ch.GetChannel()
	return &Link{
		Link: model.Link{
			Role:      c.Role,
			Net:       c.Net,
			Addr:      conn.RemoteAddr().String(),
			ChannelId: c.Id,
			PluginId:  c.PluginId,
			Online:    true,
			OnlineAt:  time.Now(),
		},
		channel: ch,
		conn:    conn,
	}
}

func newPacketLink(ch Channel, conn net.PacketConn, addr net.Addr) *Link {
	c := ch.GetChannel()
	return &Link{
		Link: model.Link{
			Role:      c.Role,
			Net:       c.Net,
			Addr:      addr.String(),
			ChannelId: c.Id,
			PluginId:  c.PluginId,
			Online:    true,
			OnlineAt:  time.Now(),
		},
		channel: ch,
		conn: &PackConn{
			PacketConn: conn,
			addr:       addr,
		},
	}
}
