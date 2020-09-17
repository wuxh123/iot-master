package dtu

import (
	"errors"
	"git.zgwit.com/iot/dtu-admin/base"
	"git.zgwit.com/iot/dtu-admin/db"
	"git.zgwit.com/iot/dtu-admin/model"
	"net"
	"time"
)

type Link struct {
	model.Link

	Rx int
	Tx int

	//设备连接
	conn net.Conn

	//发送缓存
	cache [][]byte


	lastTime time.Time
}

func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

	//TODO 发送至MQTT
}

func (l *Link) Resume() {
	for _, b := range l.cache {
		_, _ = l.Send(b)
	}
	l.cache = make([][]byte, 0)
}

func (l *Link) Send(buf []byte) (int, error) {
	//检查状态，如果关闭，则缓存
	if l.conn == nil {
		l.cache = append(l.cache, buf)
		return 0, errors.New("链接已关闭")
	}

	l.Tx += len(buf)
	l.lastTime = time.Now()

	n, e := l.conn.Write(buf)
	//TODO 没发完，继续发


	//TODO 发送至MQTT

	return n, e
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

	//TODO 发送至MQTT

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
		Link: model.Link{
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
		Link: model.Link{
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
