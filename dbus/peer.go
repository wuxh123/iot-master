package dbus

import (
	"fmt"
	"github.com/zgwit/dtu-admin/base"
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/packet"
	"sync"
	"time"
)

type Peer struct {
	baseClient

	link *dtu.Link
}

func (p *Peer) Handle(msg *packet.Packet) {
	switch msg.Type {
	case packet.TypeConnect:
		_ = p.Disconnect("duplicate register")
		_ = p.CLose()
	case packet.TypeHeartBeak:
	case packet.TypePing:
		_ = p.Send(&packet.Packet{Type: packet.TypePong})
	case packet.TypeTransfer:
		p.handleTransfer(msg)
	default:
		_ = p.SendError(fmt.Sprintf("unknown command %d", msg.Type))
	}
}

func (p *Peer) handleTransfer(msg *packet.Packet) {
	//_, _ = p.link.Send(msg.Data)
	//TODO 判断link是否为空
	//TODO 如果link断线，缓存 数据包

}

var peerKeys sync.Map // <key string, *Link>

//TODO 要检查link是否正在透传
func PreparePeer(link *dtu.Link) string {
	key := base.RandomString(20)
	//检查重复（基本不需要）
	peerKeys.Store(key, link)
	//1分钟自动清除
	time.AfterFunc(time.Minute, func() {
		peerKeys.Delete(key)
	})
	return key
}
