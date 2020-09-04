package dtu

import (
	"errors"
	"fmt"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/model"
	"log"
	"net"
	"sync"
	"time"
)

type Channel struct {
	model.Channel

	Rx int
	Tx int

	listener net.Listener

	packetConn net.PacketConn

	packetIndexes sync.Map //<Link>

	links sync.Map
}

func NewChannel(channel *model.Channel) *Channel {
	return &Channel{
		Channel: *channel,
	}
}

func (c *Channel) Open() error {
	switch c.Role {
	case "server":
		return c.Listen()
	case "client":
		return c.Dial()
	default:
		return errors.New("未知角色")
	}
}

func (c *Channel) Dial() error {
	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		c.Error = err.Error()
		return err
	}

	go c.receive(conn)

	//TODO 自动重连机制

	return err
}

func (c *Channel) Listen() error {
	var err error
	switch c.Net {
	case "tcp", "tcp4", "tcp6", "unix":
		c.listener, err = net.Listen(c.Net, c.Addr)
		if err != nil {
			c.Error = err.Error()
			return err
		}
		go c.accept()

	case "udp", "udp4", "udp6", "unixgram":
		c.packetConn, err = net.ListenPacket(c.Net, c.Addr)

		if err != nil {
			c.Error = err.Error()
			return err
		}
		go c.receivePacket()
	default:
		return errors.New("未知的网络类型")
	}
	return nil
}

func (c *Channel) Close() error {
	if c.listener != nil {
		err := c.listener.Close()
		if err != nil {
			return err
		}
		c.listener = nil

		//TODO 删除子连接
	}
	if c.packetConn != nil {
		err := c.packetConn.Close()
		if err != nil {
			return err
		}
		c.packetConn = nil
	}
	return nil
}

func (c *Channel) GetLink(id int64) (*Link, error) {
	v, ok := c.links.Load(id)
	if !ok {
		return nil, errors.New("连接不存在")
	}
	return v.(*Link), nil
}

func (c *Channel) accept() {
	for c.listener != nil {
		conn, err := c.listener.Accept()
		if err != nil {
			log.Println("accept fail:", err)
			continue
		}
		go c.receive(conn)
	}
}

func (c *Channel) receive(conn net.Conn) {
	link := newLink(c, conn)

	//未开启注册，则直接保存
	if !c.RegisterEnable {
		c.storeLink(link)
	}

	buf := make([]byte, 1024)
	for link.conn != nil {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println(e)
			break
		}
		link.onData(buf[:n])
	}

	//删除connect，或状态置空
	err := link.Close()
	if err != nil {
		log.Println(err)
	}

	if c.Role == "server" && link.Serial != "" {
		c.links.Delete(link.Id)
	} else {
		//等待5分钟，之后设为离线
		time.AfterFunc(time.Minute*5, func() {
			c.links.Delete(link.Id)
		})
	}
}

func (c *Channel) storeLink(l *Link) {
	//保存链接
	_, err := db.Engine.Insert(&l.Link)
	if err != nil {
		log.Println(err)
	}

	//根据ID保存
	c.links.Store(c.Id, l)
}

func (c *Channel) storeError(err error) error {
	c.Error = err.Error()
	_, err = db.Engine.ID(c.Id).Cols("error").Update(&c.Channel)
	return err
}

func (c *Channel) receivePacket() {
	buf := make([]byte, 1024)
	for c.packetConn != nil {
		n, addr, err := c.packetConn.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		key := addr.String()

		//找到连接，将消息发送过去
		var client *Link
		v, ok := c.packetIndexes.Load(key)
		if ok {
			client = v.(*Link)
		} else {
			client = newPacketLink(c, c.packetConn, addr)

			//根据ID保存
			if !c.RegisterEnable {
				c.storeLink(client)
			}

			//根据地址保存，收到UDP包之后，方便索引
			c.packetIndexes.Store(key, client)
		}

		client.onData(buf[:n])
	}
}
