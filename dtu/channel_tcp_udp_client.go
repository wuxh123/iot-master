package dtu

import (
	"bytes"
	"encoding/hex"
	"git.zgwit.com/iot/dtu-admin/db"
	"git.zgwit.com/iot/dtu-admin/model"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
	"log"
	"net"
	"sync"
	"time"
)

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
	err := db.DB("link").Select(q.Eq("ChannelId", c.Id), q.Eq("Role", "client")).First(&link)
	if err == storm.ErrNotFound {
		//找不到
	} else if err != nil {
		_, _ = client.Send([]byte(err.Error()))
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
			if bytes.Compare(b, buf[:n]) == 0 {
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

