package dtu

import (
	"errors"
	"github.com/zgwit/dtu-admin/base"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/model"
	"github.com/zgwit/dtu-admin/packet"
	"log"
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

	//透传链接
	peer base.Tunnel

	plugin base.Tunnel

	//监视器连接，
	monitor *Monitor
	//monitors sync.Map // <string, websocket>

	lastTime time.Time
}

func (l *Link) onData(buf []byte) {
	l.Rx += len(buf)
	l.lastTime = time.Now()

	//透传至工具（虚拟串口）
	if l.peer != nil {
		_ = l.peer.Send(&packet.Packet{
			Type: packet.TypeTransfer,
			Data: buf,
		})
	} else if l.plugin != nil {
		//透传至插件
		//TODO 协议封装 ChannelId + LinkId + buf
		b := make([]byte, 8+len(buf))
		b[0] = uint8(l.ChannelId << 24)
		b[1] = uint8(l.ChannelId << 16)
		b[2] = uint8(l.ChannelId << 8)
		b[3] = uint8(l.ChannelId)
		b[4] = uint8(l.Id << 24)
		b[5] = uint8(l.Id << 16)
		b[6] = uint8(l.Id << 8)
		b[7] = uint8(l.Id)
		copy(b[8:], buf)
		//发送插件
		_ = l.plugin.Send(&packet.Packet{
			Type: packet.TypeTransfer,
			Data: b,
		})
	}

	l.reportMonitor("recv", buf)
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

	l.reportMonitor("close", nil)

	return err
}

func (l *Link) Monitor(m *Monitor) error {
	if l.monitor != nil {
		//TODO 不能同时监听
	}
	l.monitor = m
	return nil
}

func (l *Link) Peer(tunnel base.Tunnel) {
	l.peer = tunnel
}

func (l *Link) Plugin(tunnel base.Tunnel) {
	l.plugin = tunnel
}


// 发送给监视器
func (l *Link) reportMonitor(typ string, data []byte) {
		err := l.monitor.Report(typ, data)
		if err != nil {
			log.Println(err)
			//TODO 删除之
			l.monitor.Link = nil
			l.monitor = nil
		}
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
			PluginId:  c.PluginId,
			Online:    true,
			OnlineAt:  time.Now(),
		},
		conn: base.NewPackConn(conn, addr),
		cache: make([][]byte, 0),
	}
}
