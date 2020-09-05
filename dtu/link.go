package dtu

import (
	"errors"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/model"
	"log"
	"net"
	"sync"
	"time"
)

type Link struct {
	model.Link

	Rx int
	Tx int

	//设备连接
	conn net.Conn

	//透传链接
	peer net.Conn

	//监视器连接，
	monitors sync.Map // <string, websocket>


	lastTime time.Time
}

func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

	if l.peer != nil {
		//TODO 协议封装 ChannelId + LinkId + recv + buf
		_, _ = l.peer.Write(buf)
	}

	l.reportMonitor("recv", buf)
}

func (l *Link) Send(buf []byte) (int, error) {
	l.Tx += len(buf)
	l.lastTime = time.Now()

	n, e := l.conn.Write(buf)

	if l.peer != nil {
		//TODO 协议封装
		_, _ = l.peer.Write(buf)
	}

	l.reportMonitor("send", buf)

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
	l.Online = false
	_, err = db.Engine.ID(l.Id).Cols("online").Update(&l.Link)


	l.reportMonitor("send", nil)

	return err
}

func (l *Link) Monitor(m *Monitor) {
	l.monitors.Store(m, true)
}

// 发送给监视器
func (l *Link) reportMonitor(typ string, data []byte)  {
	l.monitors.Range(func(key, value interface{}) bool {
		m := value.(*Monitor)
		err := m.Report(typ, data)
		if err != nil {
			log.Println(err)
			l.monitors.Delete(key)
		}
		return true
	})
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
