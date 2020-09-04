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

	//RemoteAddr net.Addr

	Rx int
	Tx int

	conn net.Conn

	lastTime time.Time

	channel *Channel
}

func (l *Link) checkRegister(buf []byte) error {
	n := len(buf)
	if n < l.channel.RegisterMin {
		return fmt.Errorf("register package is too short %d %s", n, string(buf[:n]))
	}
	serial := string(buf[:n])
	if l.channel.RegisterMax > 0 && l.channel.RegisterMax >= l.channel.RegisterMin && n > l.channel.RegisterMax {
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
		lnk, err := l.channel.GetLink(link.Id)
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

		link.Addr = l.conn.RemoteAddr().String()
		link.Online = true
		link.OnlineAt = time.Now()
		link.Error = ""

		_, err = db.Engine.ID(link.Id).Cols("addr", "error", "online", "online_at").Update(link)
		if err != nil {
			return err
		}
	} else {
		//插入新记录
		_, err := db.Engine.Insert(&l.Link)
		if err != nil {
			return err
		}
	}

	//保存链接
	l.channel.links.Store(link.Id, l)

	//处理剩余内容
	if l.channel.RegisterMax > 0 && n > l.channel.RegisterMax {
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

func newLink(c *Channel, conn net.Conn) *Link {
	return &Link{
		Link: model.Link{
			Addr:      conn.RemoteAddr().String(),
			ChannelId: c.Id,
			PluginId:  c.PluginId,
			Online:    true,
			OnlineAt:  time.Now(),
		},
		channel: c,
		conn:    conn,
	}
}

func newPacketLink(c *Channel, conn net.PacketConn, addr net.Addr) *Link {
	return &Link{
		Link: model.Link{
			Addr:      addr.String(),
			ChannelId: c.Id,
			PluginId:  c.PluginId,
			Online:    true,
			OnlineAt:  time.Now(),
		},
		channel: c,
		conn: &PackConn{
			PacketConn: conn,
			addr:       addr,
		},
	}
}
