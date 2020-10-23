package core

import (
	"errors"
	"fmt"
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
	"log"
	"net"
	"sync"
	"time"
)

type TcpServer struct {
	baseTunnel

	clients sync.Map
	listener net.Listener
}

func (s *TcpServer) Open() error {
	var err error
	s.listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	go s.accept()

	return nil
}

func (s *TcpServer) Close() error {
	s.clients.Range(func(key, value interface{}) bool {
		l := value.(*Link)
		_ = l.Close()
		return true
	})
	s.clients = sync.Map{}

	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			return err
		}
		s.listener = nil
	}

	return nil
}

func (s *TcpServer) GetLink(id int) (*Link, error) {
	v, ok := s.clients.Load(id)
	if !ok {
		return nil, errors.New("连接不存在")
	}
	return v.(*Link), nil
}

func (s *TcpServer) accept() {
	for s.listener != nil {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("accept fail:", err)
			continue
		}

		//开启注册了，支持多连接，否则，只支持一个链接
		if s.RegisterEnable {
			go s.receive(conn)
		} else {
			s.receive(conn)
		}
	}
}

func (s *TcpServer) receive(conn net.Conn) {
	link := newLink(s, conn)
	defer link.Close()

	buf := make([]byte, 1024)

	//第一个包作为注册包
	if s.RegisterEnable {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println(e)
			return
		}

		serial, err := s.baseTunnel.checkRegister(buf[:n])
		if err != nil {
			_ = link.Write([]byte(err.Error()))
			return
		}

		//配置序列号
		link.Serial = serial

		//查找数据库同通道，同序列号链接，更新数据库中 addr online
		var lnk models.Link
		err = db.DB("link").Select(q.Eq("TunnelId", s.ID), q.Eq("Serial", serial)).First(&link)
		if err == storm.ErrNotFound {
			//找不到
		} else if err != nil {
			_ = link.Write([]byte(err.Error()))
			log.Println(err)
			return
		} else {
			//复用连接，更新地址，状态，等
			l, _ := s.GetLink(lnk.ID)
			if l != nil {
				//如果同序号连接还在正常通讯，则关闭当前连接
				if l.conn != nil {
					_ = link.Write([]byte(fmt.Sprintf("duplicate serial %s", serial)))
					return
				}

				//复制有用的历史数据
				//link.Rx = l.Rx
				//link.Tx = l.Tx

				//复制watcher

				link.Resume()
			}

			link.ID = lnk.ID
			//link.Name = lnk.Name
		}

		//处理剩余内容
		if s.RegisterMax > 0 && n > s.RegisterMax {
			link.onData(buf[s.RegisterMax:n])
		}
	}

	//保存链接
	//s.storeLink(link)
	s.clients.Store(link.ID, link)

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
		s.clients.Delete(link.ID)
	} else {
		//有序号，等待5分钟，之后设为离线
		time.AfterFunc(time.Minute*5, func() {
			lnk, _ := s.GetLink(link.ID)
			//判断指针地址也行
			if lnk != nil && lnk.conn == nil {
				s.clients.Delete(link.ID)
			}
		})
	}
}

