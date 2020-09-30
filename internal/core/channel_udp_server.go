package core

import (
	"errors"
	"fmt"
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/internal/types"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
	"log"
	"net"
	"sync"
)

type PacketServer struct {
	baseChannel

	//处理UDP Server
	packetConn    net.PacketConn
	packetIndexes sync.Map //<Link>

}

func (c *PacketServer) Open() error {
	var err error
	c.packetConn, err = net.ListenPacket(c.Net, c.Addr)

	if err != nil {
		_ = c.storeError(err)
		return err
	}
	go c.receive()
	return nil
}

func (c *PacketServer) Close() error {

	if c.packetConn != nil {
		err := c.packetConn.Close()
		if err != nil {
			return err
		}
		c.packetConn = nil
	}
	c.clients = sync.Map{}
	return nil
}

func (c *PacketServer) GetLink(id int) (*Link, error) {
	v, ok := c.clients.Load(id)
	if !ok {
		return nil, errors.New("连接不存在")
	}
	return v.(*Link), nil
}

func (c *PacketServer) receive() {
	buf := make([]byte, 1024)
	for c.packetConn != nil {
		n, addr, err := c.packetConn.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		key := addr.String()

		//找到连接，将消息发送过去
		var link *Link
		v, ok := c.packetIndexes.Load(key)
		if ok {
			link = v.(*Link)

			//处理数据
			link.onData(buf[:n])
		} else {
			link = newPacketLink(c, c.packetConn, addr)

			//第一个包作为注册包
			if c.RegisterEnable {
				serial, err := c.baseChannel.checkRegister(buf[:n])
				if err != nil {
					_ = link.Write([]byte(err.Error()))
					return
				}

				//配置序列号
				link.Serial = serial

				//查找数据库同通道，同序列号链接，更新数据库中 addr online
				var lnk types.Link

				err = db.DB("link").Select(q.Eq("ChannelId", c.Id), q.Eq("Serial", serial)).First(&lnk)
				if err == storm.ErrNotFound {
					//找不到
				} else if err != nil {
					_ = link.Write([]byte(err.Error()))
					log.Println(err)
					return
				} else {
					//复用连接，更新地址，状态，等
					l, _ := c.GetLink(lnk.Id)
					if l != nil {
						//如果同序号连接还在正常通讯，则关闭当前连接
						if l.conn != nil {
							_ = link.Write([]byte(fmt.Sprintf("duplicate serial %s", serial)))
							return
						}

						//复制有用的历史数据
						link.Rx = l.Rx
						link.Tx = l.Tx

						//复制watcher
						link.Resume()
					}

					link.Id = lnk.Id
					link.Name = lnk.Name
				}

				//处理剩余内容
				if c.RegisterMax > 0 && n > c.RegisterMax {
					link.onData(buf[c.RegisterMax:n])
				}
			} else {
				link.onData(buf[:n])
			}

			//保存链接
			c.storeLink(link)

			//根据地址保存，收到UDP包之后，方便索引
			c.packetIndexes.Store(key, link)

			//TODO 超时自动断线，应该在一个独立的线程中检查
		}
	}
}

