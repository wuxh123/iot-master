package dtu

import (
	"errors"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/model"
	"net"
	"time"
)

type Link struct {
	model.Link

	registerChecked bool

	Rx int
	Tx int

	conn net.Conn

	lastTime time.Time
}

func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

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
		conn: &PackConn{
			PacketConn: conn,
			addr:       addr,
		},
	}
}
