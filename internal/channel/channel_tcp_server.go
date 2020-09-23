package channel

import (
	"errors"
	"fmt"
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/types"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
	"log"
	"net"
	"sync"
	"time"
)

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

		serial, err := c.baseChannel.checkRegister(buf[:n])
		if err != nil {
			_, _ = link.Send([]byte(err.Error()))
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
			_, _ = link.Send([]byte(err.Error()))
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
			link.onData(buf[c.RegisterMax:n])
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

