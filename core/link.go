package core

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"git.zgwit.com/iot/beeq/packet"
	"git.zgwit.com/zgwit/iot-admin/base"
	"git.zgwit.com/zgwit/iot-admin/models"
	"log"
	"net"
	"time"
)

type Link struct {
	models.Link

	//指向通道
	tunnel Tunnel

	//设备连接
	conn net.Conn

	//发送缓存
	cache [][]byte

	peer     base.Link
	listener base.LinkListener

	lastTime time.Time
}

func (l *Link) Listen(listener base.LinkListener) {
	l.listener = listener
}

func (l *Link) onData(buf []byte) {
	//过滤心跳
	c := l.tunnel.GetTunnel()
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
	//ln := len(buf)
	//l.Rx += ln
	//l.tunnel.GetTunnel().Rx += ln

	l.lastTime = time.Now()

	//透传
	if l.peer != nil {
		_ = l.peer.Write(buf)
		return
	}

	if l.listener != nil {
		l.listener.OnLinkData(buf)
	}

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/recv", l.TunnelId, l.ID)))
	pub.SetPayload(buf)
	hive.Publish(pub)
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

	//ln := len(buf)
	//l.Tx += ln
	//l.tunnel.GetTunnel().Tx += ln

	l.lastTime = time.Now()

	_, e := l.conn.Write(buf)
	//TODO 如果没发完，继续发

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/send", l.TunnelId, l.ID)))
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

	if l.peer != nil {
		_ = l.peer.Detach()
	}

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/event", l.TunnelId, l.ID)))
	pub.SetPayload([]byte("close"))
	Hive().Publish(pub)

	return err
}

func (l *Link) Attach(link base.Link) error {
	//check peer
	l.peer = link
	return nil
}
func (l *Link) Detach() error {
	//check peer
	if l.peer != nil {
		_ = l.peer.Detach()
		l.peer = nil
	}
	return nil
}

func newLink(ch Tunnel, conn net.Conn) *Link {
	c := ch.GetTunnel()
	return &Link{
		Link: models.Link{
			TunnelId:  c.ID,
			//ProjectId: c.ProjectId,
			Active:    true,
		},
		tunnel: ch,
		conn:   conn,
		cache:  make([][]byte, 0),
	}
}

func newPacketLink(ch Tunnel, conn net.PacketConn, addr net.Addr) *Link {
	c := ch.GetTunnel()
	return &Link{
		Link: models.Link{
			TunnelId:  c.ID,
			//ProjectId: c.ProjectId,
			Active:    true,
		},
		tunnel: ch,
		conn:   base.NewPackConn(conn, addr),
		cache:  make([][]byte, 0),
	}
}
