package core

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"git.zgwit.com/iot/beeq/packet"
	"git.zgwit.com/zgwit/iot-admin/interfaces"
	"git.zgwit.com/zgwit/iot-admin/internal/base"
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/internal/types"
	"log"
	"net"
	"time"
)

type Link struct {
	types.LinkExt

	//指向通道
	channel Channel

	//设备连接
	conn net.Conn

	//发送缓存
	cache [][]byte

	peer *Peer

	lastTime time.Time

	listener interfaces.LinkerListener
}

func (l *Link) Listen(listener interfaces.LinkerListener) {
	l.listener = listener
}

func (l *Link) onData(buf []byte) {
	//过滤心跳
	c := l.channel.GetChannel()
	if c.HeartBeatEnable && time.Now().Sub(l.lastTime) > time.Second*time.Duration(c.HeartBeatInterval) {
		var b []byte
		if c.HeartBeatIsHex {
			var e error
			b, e = hex.DecodeString(c.HeartBeatContent)
			if e != nil {
				log.Println(e)
			}
		} else {
			b = []byte(c.HeartBeatContent)
		}
		if bytes.Compare(b, buf) == 0 {
			return
		}
	}

	//计数
	ln := len(buf)
	l.Rx += ln
	l.channel.GetChannel().Rx += ln

	l.lastTime = time.Now()

	//透传
	if l.peer != nil {
		_ = l.peer.Send(buf)
		return
	}

	//监听
	if l.listener != nil {
		l.listener.OnLinkerData(buf)
	}

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/recv", l.ChannelId, l.Id)))
	pub.SetPayload(buf)
	Hive().Publish(pub)
}

func (l *Link) Resume() {
	for _, b := range l.cache {
		_ = l.Write(b)
	}
	l.cache = make([][]byte, 0)
}

func (l *Link) Write(buf []byte) error {
	//检查状态，如果关闭，则缓存
	if l.conn == nil {
		l.cache = append(l.cache, buf)
		return errors.New("链接已关闭")
	}

	ln := len(buf)
	l.Tx += ln
	l.channel.GetChannel().Tx += ln

	l.lastTime = time.Now()

	_, e := l.conn.Write(buf)
	//TODO 如果没发完，继续发

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/send", l.ChannelId, l.Id)))
	pub.SetPayload(buf)
	Hive().Publish(pub)

	return e
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

	//监听关闭
	if l.listener != nil {
		l.listener.OnLinkerClose()
	}

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/event", l.ChannelId, l.Id)))
	pub.SetPayload([]byte("close"))
	Hive().Publish(pub)

	return err
}

func (l *Link) storeError(err error) error {
	l.Error = err.Error()
	//_, err = db.Engine.ID(l.Id).Cols("error").Update(&l.Link)
	return db.DB("link").UpdateField(&l.Link, "error", l.Error)
}

func newLink(ch Channel, conn net.Conn) *Link {
	c := ch.GetChannel()
	return &Link{
		LinkExt: types.LinkExt{
			Link: types.Link{
				Net:       c.Net,
				Addr:      conn.RemoteAddr().String(),
				ChannelId: c.Id,
			},
			Online: true,
		},
		channel: ch,
		conn:    conn,
		cache:   make([][]byte, 0),
	}
}

func newPacketLink(ch Channel, conn net.PacketConn, addr net.Addr) *Link {
	c := ch.GetChannel()
	return &Link{
		LinkExt: types.LinkExt{
			Link: types.Link{
				Net:       c.Net,
				Addr:      addr.String(),
				ChannelId: c.Id,
			},
			Online: true,
		},
		channel: ch,
		conn:    base.NewPackConn(conn, addr),
		cache:   make([][]byte, 0),
	}
}
