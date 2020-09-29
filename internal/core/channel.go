package core

import (
	"fmt"
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/types"
	"log"
	"regexp"
	"sync"
)

type Channel interface {
	Open() error
	Close() error
	GetLink(id int) (*Link, error)
	GetChannel() *types.Channel
}

func NewChannel(channel *types.Channel) (Channel, error) {
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
	types.Channel

	clients sync.Map

	Rx int `json:"rx"`
	Tx int `json:"tx"`
}

func (c *baseChannel) GetChannel() *types.Channel {
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
