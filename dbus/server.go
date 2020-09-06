package dbus

import (
	"errors"
	"github.com/zgwit/dtu-admin/base"
	"github.com/zgwit/dtu-admin/packet"
	"log"
	"net"
	"strings"
)

type Server struct {
	Net      string
	Addr     string
	listener net.Listener
}

func (s *Server) Listen() error {
	var err error
	s.listener, err = net.Listen(s.Net, s.Addr)
	if err != nil {
		return err
	}
	go s.accept()
	return nil
}

func (s *Server) Close() error {
	if s.listener != nil {
		return nil
	}
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.listener = nil
	return nil
}

func (s *Server) accept() {
	for s.listener != nil {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
			//TODO 判断监听异常应该退出
		}
		go s.receive(conn)
	}
}

func (s *Server) receive(conn net.Conn) {

	var parser packet.Parser
	buf := make([]byte, 1024)

	//接收第一个包，作为类型校验
	n, e := conn.Read(buf)
	if e != nil {
		log.Println(e)
		return
	}
	packs := parser.Parse(buf[:n])
	if len(packs) == 0 {
		_ = conn.Close()
		return
	}
	//根据第一个包创建客户羰
	c, e := s.createTunnel(packs[0])
	if e != nil {
		_, _ = conn.Write([]byte(e.Error()))
		_ = conn.Close()
		return
	}
	//处理第一次接收的剩余包（网络拥堵时，可能会发生）
	for _, pack := range packs[1:] {
		c.Handle(pack)
	}

	for {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println(e)
			break
		}
		packs := parser.Parse(buf[:n])
		for _, pack := range packs {
			c.Handle(pack)
		}
	}

	//关闭
	_ = c.CLose()
}

func (s *Server) createTunnel(p *packet.Packet) (base.Tunnel, error) {
	if p.Type != packet.TypeConnect {
		//告诉客户端，第一个包必须是注册包
		return nil, errors.New("first packet must be connect")
	}

	register := string(p.Data)
	rs := strings.Split(register, ":")
	switch rs[0] {
	case "peer":
		//TODO 解析 peer:key

	case "plugin":
		//TODO 解析 plugin:key:secret

	}
	return nil, errors.New("未支持的类型")
}

func Start(addr string) error {
	s := &Server{
		Net:  "tcp",
		Addr: addr,
	}
	return s.Listen()
}
