package tunnel

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"git.zgwit.com/iot/beeq/packet"
	"iot-master/dbus"
	"iot-master/model"
	"iot-master/protocol"
	"iot-master/types"
	"log"
	"net"
	"sync"
	"time"
)

type link struct {
	model.Link

	//指向通道
	tunnel types.Tunnel

	//设备连接
	conn net.Conn

	//发送缓存
	cache [][]byte

	peer types.OnDataFunc

	listener types.LinkListener

	lastTime time.Time

	//协议
	adapter protocol.Adapter

	//项目
	project types.Project

	//设备，以从站号为KEY
	devices sync.Map //<slave, Device>
}

func (l *link) onData(buf []byte) {
	//过滤心跳
	c := l.tunnel.GetModel()
	if c.HeartBeat.Enable && time.Now().Sub(l.lastTime) > time.Second*time.Duration(c.HeartBeat.Interval) {
		var b []byte
		if c.HeartBeat.IsHex {
			var e error
			b, e = hex.DecodeString(c.HeartBeat.Content)
			if e != nil {
				log.Println(e)
			}
		} else {
			b = []byte(c.HeartBeat.Content)
		}
		if bytes.Compare(b, buf) == 0 {
			return
		}
	}

	//计数
	//ln := len(buf)
	//l.Rx += ln
	//l.tunnel.GetModel().Rx += ln

	l.lastTime = time.Now()

	//透传
	if l.peer != nil {
		l.peer(buf)
		return
	}

	//响应数据
	if l.listener != nil {
		l.listener.OnData(buf)
	}

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/recv", l.TunnelId, l.Id)))
	pub.SetPayload(buf)
	dbus.Hive().Publish(pub)
}

func (l *link) Resume() {
	for _, b := range l.cache {
		_ = l.Write(b)
	}
	l.cache = make([][]byte, 0)
}

func (l *link) Write(buf []byte) error {
	//检查状态，如果关闭，则缓存
	if l.conn == nil {
		l.cache = append(l.cache, buf)
		return errors.New("链接已关闭")
	}

	//ln := len(buf)
	//l.Tx += ln
	//l.tunnel.GetModel().Tx += ln

	l.lastTime = time.Now()

	_, e := l.conn.Write(buf)
	//TODO 如果没发完，继续发

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/send", l.TunnelId, l.Id)))
	pub.SetPayload(buf)
	dbus.Hive().Publish(pub)

	return e
}

func (l *link) Close() error {
	if l.conn == nil {
		return errors.New("连接已经关闭")
	}
	err := l.conn.Close()
	l.conn = nil
	if err != nil {
		return err
	}

	l.peer = nil

	//发送至MQTT
	pub := packet.PUBLISH.NewMessage().(*packet.Publish)
	pub.SetTopic([]byte(fmt.Sprintf("/link/%d/%d/event", l.TunnelId, l.Id)))
	pub.SetPayload([]byte("close"))
	dbus.Hive().Publish(pub)

	return err
}

func (l *link) Attach(listener types.OnDataFunc) error {
	//check peer
	l.peer = listener
	return nil
}
func (l *link) Detach() error {
	//check peer
	if l.peer != nil {
		l.peer = nil
	}
	return nil
}

func (l *link) Listen(listener types.LinkListener) {
	l.listener = listener
}

func newLink(t types.Tunnel, conn net.Conn) *link {
	c := t.GetModel()
	return &link{
		Link: model.Link{
			TunnelId: c.Id,
			Active:   true,
		},
		tunnel: t,
		conn:   conn,
		cache:  make([][]byte, 0),
	}
}

func newPacketLink(ch types.Tunnel, conn net.PacketConn, addr net.Addr) *link {
	c := ch.GetModel()
	return &link{
		Link: model.Link{
			TunnelId: c.Id,
			Active:   true,
		},
		tunnel: ch,
		conn:   NewPackConn(conn, addr),
		cache:  make([][]byte, 0),
	}
}
