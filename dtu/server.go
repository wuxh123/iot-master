package dtu

import (
	"errors"
	"log"
	"net"
	"regexp"
	"sync"
)

type Server struct {
	Net  string
	Addr string

	listener net.Listener
	clients  sync.Map

	increment int

	RegisterPack struct {
		Enable    bool
		MinLength int
		MaxLength int
		Regex     string
	}

	HeartBeatPack struct {
		Enable  bool
		Content []byte
	}
}

func NewServer(net string, addr string) *Server {
	return &Server{
		Net:  net,
		Addr: addr,
	}
}

func (s *Server) Open() error {
	var err error
	s.listener, err = net.Listen(s.Net, s.Addr)
	if err != nil {
		return err
	}
	go s.accept()
	return nil
}

func (s *Server) Close() error {
	if s.listener == nil {
		return errors.New("server closed")
	}
	err := s.listener.Close()
	s.clients.Range(func(k, v interface{}) bool {
		client := v.(*Client)
		_ = client.Close()
		return true
	})
	s.listener = nil
	return err
}

func (s *Server) accept() {
	for s.listener != nil {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("accept fail:", err)
			continue
		}

		go s.receive(conn)
	}
}

func (s *Server) receive(conn net.Conn) {

	if !s.RegisterPack.Enable {
		//匿名链接
		client := NewClient(s.Net, s.Addr)
		s.increment++
		s.clients.Store(s.increment, client)

		client.conn = conn
		client.receive()
		return
	}

	//接收注册包
	sn := make([]byte, s.RegisterPack.MaxLength)
	n, e := conn.Read(sn)
	if e != nil {
		log.Println("read error", e)
		return
	}
	if n < s.RegisterPack.MinLength {
		log.Println("register package is too short", n, string(sn[:n]))
		_ = conn.Close()
		return
	}
	serial := string(sn[:n])

	// 正则表达式判断合法性
	if s.RegisterPack.Regex != "" {
		reg := regexp.MustCompile(`^` + s.RegisterPack.Regex + `$`)
		match := reg.MatchString(serial)
		if !match {
			log.Println("register package format error", serial)
			_ = conn.Close()
			return
		}
	}

	//根据SN查找缓存的链接
	var client *Client
	v, ok := s.clients.Load(serial)
	if ok {
		client = v.(*Client)
	} else {
		client = NewClient(s.Net, s.Addr)
		s.clients.Store(serial, client)
		//client.Id = mh.Id
	}
	client.conn = conn
	client.receive()
}
