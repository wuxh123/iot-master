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
	client   *Link    //作为客户端的连接
	clients  sync.Map //接入的客户端

	//处理UDP Server
	packetConn    net.PacketConn
	packetIndexes sync.Map //<Link>

}

func NewChannel(channel *model.Channel) *Channel {
	return &Channel{
		Channel: *channel,
	}
}

func (c *Channel) Open() error {
	switch c.Role {
	case "server":
		if c.listener != nil {
			return errors.New("已经打开监听了")
		}
		return c.Listen()
	case "client":
		if c.client != nil {
			return errors.New("已经连接了")
		}
		return c.Dial()
	default:
		return errors.New("未知角色")
	}
}

func (c *Channel) Dial() error {
	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		_ = c.storeError(err)
		return err
	}

	go c.receiveClient(conn)

	return nil
}

func (c *Channel) Listen() error {
	var err error
	switch c.Net {
	case "tcp", "tcp4", "tcp6", "unix":
		c.listener, err = net.Listen(c.Net, c.Addr)
		if err != nil {
			_ = c.storeError(err)
			return err
		}
		go c.accept()

	case "udp", "udp4", "udp6", "unixgram":
		c.packetConn, err = net.ListenPacket(c.Net, c.Addr)

		if err != nil {
			_ = c.storeError(err)
			return err
		}
		go c.receivePacket()
	default:
		return errors.New("未知的网络类型")
	}
	return nil
}

func (c *Channel) Close() error {
	c.clients.Range(func(key, value interface{}) bool {
		l := value.(*Link)
		_ = l.Close()
		return true
	})
	c.clients = sync.Map{}

	if c.client != nil {
		err := c.client.Close()
		if err != nil {
			return err
		}
		c.client = nil
	}

	if c.listener != nil {
		err := c.listener.Close()
		if err != nil {
			return err
		}
		c.listener = nil
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
	v, ok := c.clients.Load(id)
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
		go c.receiveServer(conn)
	}
}

func (c *Channel) receiveClient(conn net.Conn) {
	c.client = newLink(c, conn)

	buf := make([]byte, 1024)
	for c.client != nil && c.client.conn != nil {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println(e)
			break
		}
		c.client.onData(buf[:n])
	}

	//关闭了通道，就再不管了
	if c.client == nil {
		return
	}

	if c.client.conn != nil {
		err := c.client.Close()
		if err != nil {
			log.Println(err)
		}
	}

	//重连，要判断是否是关闭状态，计算重连次数
	_ = c.Dial()
}

func (c *Channel) receiveServer(conn net.Conn) {
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

	//无序号，直接删除
	if link.Serial != "" {
		c.clients.Delete(link.Id)
	} else {
		//有序号，等待5分钟，之后设为离线
		time.AfterFunc(time.Minute*5, func() {
			c.clients.Delete(link.Id)
		})
	}
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

func (c *Channel) storeLink(l *Link) {
	//保存链接
	_, err := db.Engine.Insert(&l.Link)
	if err != nil {
		log.Println(err)
	}

	//根据ID保存
	c.clients.Store(c.Id, l)
}

func (c *Channel) storeError(err error) error {
	c.Error = err.Error()
	_, err = db.Engine.ID(c.Id).Cols("error").Update(&c.Channel)
	return err
}
