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

type Channel interface {
	Open() error
	Close() error
	GetLink(id int64) (*Link, error)
	GetChannel() *model.Channel
	StoreLink(link *Link)
}

func NewChannel(channel *model.Channel) (Channel, error) {
	if channel.Role == "client" {
		return &Client{
			baseChannel: baseChannel{
				Channel: *channel,
			},
		}, nil
	} else if channel.Role == "server" {
		switch channel.Net {
		case "tcp", "tcp4", "tcp6", "unix":
			return &Server{
				baseChannel: baseChannel{
					Channel: *channel,
				},
			}, nil
		case "udp", "udp4", "udp6", "unixgram":
			return &PacketServer{
				baseChannel: baseChannel{
					Channel: *channel,
				},
			}, nil
		default:
			return nil, fmt.Errorf("未知的网络类型 %s", channel.Net)
		}
	} else {
		return nil, fmt.Errorf("未知的角色 %s", channel.Role)
	}
}

type baseChannel struct {
	model.Channel

	clients sync.Map

	Rx int `json:"rx"`
	Tx int `json:"tx"`
}

func (c *baseChannel) GetChannel() *model.Channel {
	return &c.Channel
}

func (c *baseChannel) StoreLink(l *Link) {
	//保存链接
	if l.Id > 0 {
		_, err := db.Engine.ID(l.Id).Cols("addr", "error", "online", "online_at").Update(&l.Link)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := db.Engine.Insert(&l.Link)
		if err != nil {
			log.Println(err)
		}
	}

	//根据ID保存
	c.clients.Store(l.Id, l)
}

func (c *baseChannel) storeError(err error) error {
	c.Error = err.Error()
	_, err = db.Engine.ID(c.Id).Cols("error").Update(&c.Channel)
	return err
}

type Client struct {
	baseChannel

	client *Link //作为客户端的连接
}

func (c *Client) Open() error {
	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		_ = c.storeError(err)
		return err
	}

	go c.receiveClient(conn)

	return nil
}

func (c *Client) Close() error {

	c.clients = sync.Map{}

	if c.client != nil {
		err := c.client.Close()
		if err != nil {
			return err
		}
		c.client = nil
	}

	return nil
}

func (c *Client) GetLink(id int64) (*Link, error) {
	return c.client, nil
}

func (c *Client) receiveClient(conn net.Conn) {
	c.client = newLink(c, conn)

	var link model.Link
	has, err := db.Engine.Where("channel_id=?", c.Id).And("role=?", "client").Get(&link)
	if err != nil {
		log.Println(err)
		return
	}
	if has {
		//复用连接，更新地址，状态，等
		link.Addr = conn.RemoteAddr().String()
		link.Online = true
		link.OnlineAt = time.Now()
		link.Error = ""
		_, err = db.Engine.ID(link.Id).Cols("addr", "error", "online", "online_at").Update(link)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		_, err := db.Engine.Insert(&c.client.Link)
		if err != nil {
			log.Println(err)
		}
	}

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
	_ = c.Open()
}

type Server struct {
	baseChannel

	listener net.Listener
}

func (c *Server) Open() error {
	var err error
	c.listener, err = net.Listen(c.Net, c.Addr)
	if err != nil {
		_ = c.storeError(err)
		return err
	}
	go c.accept()

	return nil
}

func (c *Server) Close() error {
	c.clients.Range(func(key, value interface{}) bool {
		l := value.(*Link)
		_ = l.Close()
		return true
	})
	c.clients = sync.Map{}

	if c.listener != nil {
		err := c.listener.Close()
		if err != nil {
			return err
		}
		c.listener = nil
	}

	return nil
}

func (c *Server) GetLink(id int64) (*Link, error) {
	v, ok := c.clients.Load(id)
	if !ok {
		return nil, errors.New("连接不存在")
	}
	return v.(*Link), nil
}

func (c *Server) accept() {
	for c.listener != nil {
		conn, err := c.listener.Accept()
		if err != nil {
			log.Println("accept fail:", err)
			continue
		}
		go c.receive(conn)
	}
}

func (c *Server) receive(conn net.Conn) {
	link := newLink(c, conn)

	//未开启注册，则直接保存
	if !c.RegisterEnable {
		c.StoreLink(link)
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

func (c *PacketServer) GetLink(id int64) (*Link, error) {
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
		var client *Link
		v, ok := c.packetIndexes.Load(key)
		if ok {
			client = v.(*Link)
		} else {
			client = newPacketLink(c, c.packetConn, addr)

			//根据ID保存
			if !c.RegisterEnable {
				c.StoreLink(client)
			}

			//根据地址保存，收到UDP包之后，方便索引
			c.packetIndexes.Store(key, client)
		}

		client.onData(buf[:n])
	}
}
