package core

import (
	"errors"
	"fmt"
	"git.zgwit.com/iot/beeq/packet"
	"git.zgwit.com/zgwit/iot-admin/interfaces"
	"git.zgwit.com/zgwit/iot-admin/internal/base"
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/internal/types"
	"net"
	"time"
)

type Link struct {
	types.Link

	Rx int
	Tx int

	//设备连接
	conn net.Conn

	//发送缓存
	cache [][]byte


	lastTime time.Time

	listener interfaces.LinkerListener
}

func (l *Link) Listen (listener interfaces.LinkerListener) {
	l.listener = listener
}


func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

	//TODO 透传

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

	l.Tx += len(buf)
	l.lastTime = time.Now()

	_, e := l.conn.Write(buf)
	//TODO 没发完，继续发

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
		Link: types.Link{
			Role:      c.Role,
			Net:       c.Net,
			Addr:      conn.RemoteAddr().String(),
			ChannelId: c.Id,
			//PluginId:  c.PluginId,
			Online:    true,
			OnlineAt:  time.Now(),
		},
		conn:  conn,
		cache: make([][]byte, 0),
	}
}

func newPacketLink(ch Channel, conn net.PacketConn, addr net.Addr) *Link {
	c := ch.GetChannel()
	return &Link{
		Link: types.Link{
			Role:      c.Role,
			Net:       c.Net,
			Addr:      addr.String(),
			ChannelId: c.Id,
			//PluginId:  c.PluginId,
			Online:    true,
			OnlineAt:  time.Now(),
		},
		conn:  base.NewPackConn(conn, addr),
		cache: make([][]byte, 0),
	}
}
