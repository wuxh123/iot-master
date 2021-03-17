package tunnel

import (
	"errors"
	"fmt"
	"iot-master/db"
	"iot-master/model"
	"log"
	"net"
	"sync"
)

type PacketServer struct {
	tunnel

	//处理UDP Server
	packetConn    net.PacketConn
	packetIndexes sync.Map //<Link>

	clients sync.Map
}

func (c *PacketServer) Open() error {
	var err error
	c.packetConn, err = net.ListenPacket("udp", c.Addr)

	if err != nil {
		//_ = c.storeError(err)
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
	c.packetIndexes = sync.Map{}
	return nil
}

func (c *PacketServer) GetLink(id int) (Link, error) {
	return c.getLink(id)
}

func (c *PacketServer) getLink(id int) (*link, error) {
	v, ok := c.clients.Load(id)
	if !ok {
		return nil, errors.New("连接不存在")
	}
	return v.(*link), nil
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
		var lnk *link
		v, ok := c.packetIndexes.Load(key)
		if ok {
			lnk = v.(*link)

			//处理数据
			lnk.onData(buf[:n])
		} else {
			lnk = newPacketLink(c, c.packetConn, addr)

			//第一个包作为注册包
			if c.RegisterEnable {
				serial, err := c.tunnel.checkRegister(buf[:n])
				if err != nil {
					_ = lnk.Write([]byte(err.Error()))
					return
				}

				//配置序列号
				lnk.Serial = serial

				//查找数据库同通道，同序列号链接，更新数据库中 addr online
				var l model.Link

				has, err := db.Engine.Where("tunnel_id", c.Id).And("serial", serial).Get(&l)
				if !has {
					//找不到
				} else if err != nil {
					_ = lnk.Write([]byte(err.Error()))
					log.Println(err)
					return
				} else {
					//复用连接，更新地址，状态，等
					ll, _ := c.getLink(lnk.Id)
					if ll != nil {
						//如果同序号连接还在正常通讯，则关闭当前连接
						if ll.conn != nil {
							_ = lnk.Write([]byte(fmt.Sprintf("duplicate serial %s", serial)))
							return
						}

						//复制有用的历史数据
						//lnk.Rx = l.Rx
						//lnk.Tx = l.Tx

						//复制watcher
						lnk.Resume()
					}

					lnk.Id = l.Id
					//lnk.Name = lnk.Name
				}

				//处理剩余内容
				if c.RegisterMax > 0 && n > c.RegisterMax {
					lnk.onData(buf[c.RegisterMax:n])
				}
			} else {
				lnk.onData(buf[:n])
			}

			//保存链接
			//c.storeLink(lnk)
			c.clients.Store(lnk.Id, lnk)

			//根据地址保存，收到UDP包之后，方便索引
			c.packetIndexes.Store(key, lnk)

			//TODO 超时自动断线，应该在一个独立的线程中检查
		}
	}
}
