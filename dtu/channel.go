package dtu

import (
	"errors"
	"fmt"
	"github.com/zgwit/dtu-admin/types"
	"log"
	"net"
	"sync"
)


type Channel struct {
	types.Channel

	Error string

	listener      net.Listener
	packetConn    net.PacketConn
	packetIndexes sync.Map

	links sync.Map

	//自增ID
	increment int64
}

func NewChannel(channel *types.Channel) *Channel {
	return &Channel{
		Channel: *channel,
	}
}

func (c *Channel) Open() error {
	if c.IsServer {
		return c.Listen()
	} else {
		return c.Dial()
	}
}

func (c *Channel) Dial() error {
	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		log.Println(err)
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
			return err
		}
		go c.accept()

	case "udp", "udp4", "udp6", "unixgram":
		c.packetConn, err = net.ListenPacket(c.Net, c.Addr)

		if err != nil {
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
	client := newConnection(conn)
	client.channel = c
	c.storeConnection(client)

	buf := make([]byte, 1024)
	for client.conn != nil {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println(e)
			break
		}
		client.onData(buf[:n])
	}

	//TODO 删除connect，或状态置空

}

func (c *Channel) storeConnection(conn *Link)  {
	c.increment++
	conn.ID = c.increment

	//根据ID保存
	c.links.Store(c.ID, conn)
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
			client = newPacketConnection(c.packetConn, addr)
			client.channel = c

			//根据ID保存
			c.storeConnection(client)

			//根据地址保存，收到UDP包之后，方便索引
			c.packetIndexes.Store(key, client)
		}

		client.onData(buf[:n])
	}
}
