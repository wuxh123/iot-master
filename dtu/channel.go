package dtu

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/model"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
	"log"
	"net"
	"regexp"
	"sync"
	"time"
)

type Channel interface {
	Open() error
	Close() error
	GetLink(id int) (*Link, error)
	GetChannel() *model.Channel
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

func (c *baseChannel) storeLink(l *Link) {
	//保存链接
	if l.Id > 0 {
		err := db.DB("link").Update(&l.Link)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := db.DB("link").Save(&l.Link)
		if err != nil {
			log.Println(err)
		}
	}

	//根据ID保存
	c.clients.Store(l.Id, l)
}

func (c *baseChannel) storeError(err error) error {
	c.Error = err.Error()
	return db.DB("channel").UpdateField(&c.Channel, "error", c.Error)
}

func (c *baseChannel) checkRegister(buf []byte) (string, error) {
	n := len(buf)
	if n < c.RegisterMin {
		return "", fmt.Errorf("register package is too short %d %s", n, string(buf[:n]))
	}
	serial := string(buf[:n])
	if c.RegisterMax > 0 && c.RegisterMax >= c.RegisterMin && n > c.RegisterMax {
		serial = string(buf[:c.RegisterMax])
	}

	// 正则表达式判断合法性
	if c.RegisterRegex != "" {
		reg := regexp.MustCompile(`^` + c.RegisterRegex + `$`)
		match := reg.MatchString(serial)
		if !match {
			return "", fmt.Errorf("register package format error %s", serial)
		}
	}

	return serial, nil
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

	go c.receive(conn)

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

func (c *Client) GetLink(id int) (*Link, error) {
	return c.client, nil
}

func (c *Client) receive(conn net.Conn) {
	//c.client = newLink(c, conn)
	client := newLink(c, conn)
	if c.client != nil {
		//重发数据
		client.cache = c.client.cache
	}
	c.client = client

	var link model.Link
	err := db.DB("link").Select(q.Eq("channel_id", c.Id), q.Eq("role", "client")).First(&link)
	if err != storm.ErrNotFound {
		//找不到
	} else if err != nil {
		log.Println(err)
		return
	} else {
		//复用连接，更新地址，状态，等
		c.client.Id = link.Id
	}

	c.storeLink(c.client)

	buf := make([]byte, 1024)
	for c.client != nil && c.client.conn != nil {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println(e)
			break
		}

		//过滤心跳包
		if c.HeartBeatEnable && time.Now().Sub(c.client.lastTime) > time.Second*time.Duration(c.HeartBeatInterval) {
			var b []byte
			if c.HeartBeatIsHex {
				var e error
				b, e = hex.DecodeString(c.HeartBeatContent)
				if e != nil {
					log.Println(e)
				}
			} else {
				b = []byte(c.HeartBeatContent)
			}
			if bytes.Compare(b, buf) == 0 {
				continue
			}
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

func (c *Server) GetLink(id int) (*Link, error) {
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
	defer link.Close()

	buf := make([]byte, 1024)

	//第一个包作为注册包
	if c.RegisterEnable {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println(e)
			return
		}

		serial, err := c.baseChannel.checkRegister(buf)
		if err != nil {
			_, _ = link.Send([]byte(err.Error()))
			return
		}

		//配置序列号
		link.Serial = serial

		//查找数据库同通道，同序列号链接，更新数据库中 addr online
		var lnk model.Link
		err = db.DB("link").Select(q.Eq("channel_id", c.Id), q.Eq("serial", serial)).First(&link)
		if err != storm.ErrNotFound {
			//找不到
		} else if err != nil {
			_, _ = link.Send([]byte("数据库异常"))
			log.Println(err)
			return
		} else {
			//复用连接，更新地址，状态，等
			l, _ := c.GetLink(lnk.Id)
			if l != nil {
				//如果同序号连接还在正常通讯，则关闭当前连接
				if l.conn != nil {
					_, _ = link.Send([]byte(fmt.Sprintf("duplicate serial %s", serial)))
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
			link.onData(buf[c.RegisterMax:])
		}
	}

	//保存链接
	c.storeLink(link)

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
			lnk, _ := c.GetLink(link.Id)
			//判断指针地址也行
			if lnk != nil && lnk.conn == nil {
				c.clients.Delete(link.Id)
			}
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

			//过滤心跳包
			if c.HeartBeatEnable && time.Now().Sub(link.lastTime) > time.Second*time.Duration(c.HeartBeatInterval) {
				var b []byte
				if c.HeartBeatIsHex {
					var e error
					b, e = hex.DecodeString(c.HeartBeatContent)
					if e != nil {
						log.Println(e)
					}
				} else {
					b = []byte(c.HeartBeatContent)
				}
				if bytes.Compare(b, buf) == 0 {
					continue
				}
			}

			//处理数据
			link.onData(buf[:n])
		} else {
			link = newPacketLink(c, c.packetConn, addr)

			//第一个包作为注册包
			if c.RegisterEnable {
				serial, err := c.baseChannel.checkRegister(buf)
				if err != nil {
					_, _ = link.Send([]byte(err.Error()))
					return
				}

				//配置序列号
				link.Serial = serial

				//查找数据库同通道，同序列号链接，更新数据库中 addr online
				var lnk model.Link

				err = db.DB("link").Select(q.Eq("channel_id", c.Id), q.Eq("serial", serial)).First(&link)
				if err != storm.ErrNotFound {
					//找不到
				} else if err != nil {
					_, _ = link.Send([]byte("数据库异常"))
					log.Println(err)
					return
				} else {
					l, _ := c.GetLink(lnk.Id)
					if l != nil {
						//如果同序号连接还在正常通讯，则关闭当前连接
						if l.conn != nil {
							_, _ = link.Send([]byte(fmt.Sprintf("duplicate serial %s", serial)))
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
					link.onData(buf[c.RegisterMax:])
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
