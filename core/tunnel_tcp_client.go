package core

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"log"
	"net"
)

type TcpClient struct {
	baseTunnel

	link *Link //作为客户端的连接
}

func (c *TcpClient) Open() error {
	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		//_ = c.storeError(err)
		return err
	}

	go c.receive(conn)

	return nil
}

func (c *TcpClient) Close() error {

	if c.link != nil {
		err := c.link.Close()
		if err != nil {
			return err
		}
		c.link = nil
	}

	return nil
}

func (c *TcpClient) GetLink(id int64) (*Link, error) {
	return c.link, nil
}

func (c *TcpClient) receive(conn net.Conn) {
	//复用地连接
	if c.link == nil {
		link := newLink(c, conn)

		var lnk models.Link
		has, err := db.Engine.Where("tunnel_id", c.Id).Get(&lnk)
		//err := db.DB("link").Select(q.Eq("ChannelId", c.Id)).First(&lnk)
		if !has {
			//找不到，新建
		} else if err != nil {
			_ = link.Write([]byte(err.Error()))
			log.Println(err)
			return
		} else {
			//复用连接，更新地址，状态，等
			c.link.Link = lnk
			//c.link.Id = lnk.Id
		}

		c.link = link
	}

	buf := make([]byte, 1024)
	for c.link != nil && c.link.conn != nil {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println(e)
			break
		}

		c.link.onData(buf[:n])
	}

	//关闭了通道，就再不管了
	if c.link == nil {
		return
	}

	if c.link.conn != nil {
		err := c.link.Close()
		if err != nil {
			log.Println(err)
		}
	}

	//重连，要判断是否是关闭状态，计算重连次数
	_ = c.Open()
}
