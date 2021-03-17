package tunnel

import (
	"errors"
	"fmt"
	"iot-master/model"
	"regexp"
	"sync"
)

var tunnels sync.Map


//通道
type Tunnel interface {
	Open() error
	Close() error
	GetTunnel() *model.Tunnel
	GetLink(id int) (Link, error)
}


type tunnel struct {
	model.Tunnel
	//model.ProjectAdapter
}

func (t *tunnel) GetTunnel() *model.Tunnel {
	return &t.Tunnel
}

func (t *tunnel) checkRegister(buf []byte) (string, error) {
	n := len(buf)
	if n < t.RegisterMin {
		return "", fmt.Errorf("register package is too short %d %s", n, string(buf[:n]))
	}
	serial := string(buf[:n])
	if t.RegisterMax > 0 && t.RegisterMax >= t.RegisterMin && n > t.RegisterMax {
		serial = string(buf[:t.RegisterMax])
	}

	// 正则表达式判断合法性
	if t.RegisterRegex != "" {
		reg := regexp.MustCompile(`^` + t.RegisterRegex + `$`)
		match := reg.MatchString(serial)
		if !match {
			return "", fmt.Errorf("register package format error %s", serial)
		}
	}

	return serial, nil
}

func NewTunnel(t *model.Tunnel) (Tunnel, error) {
	switch t.Type {
	case "tcp-server":
		return &TcpServer{
			tunnel: tunnel{
				Tunnel: *t,
			},
		}, nil
	case "tcp-client":
		return &TcpUdpClient{
			tunnel: tunnel{
				Tunnel: *t,
			},
			Net: "tcp",
		}, nil
	case "udp-server":
		return &PacketServer{
			tunnel: tunnel{
				Tunnel: *t,
			},
		}, nil
	case "udp-client":
		return &TcpUdpClient{
			tunnel: tunnel{
				Tunnel: *t,
			},
			Net: "udp",
		}, nil
	case "serial":
	default:
	}
	return nil, fmt.Errorf("未知的网络类型 %s", t.Type)
}


func StartTunnel(c *model.Tunnel) (Tunnel, error) {
	//log.Println("Start core", c)
	tunnel, err := NewTunnel(c)
	if err != nil {
		return nil, err
	}
	err = tunnel.Open()
	if err != nil {
		return nil, err
	}
	tunnels.Store(c.Id, tunnel)
	return tunnel, err
}

func DeleteTunnel(id int) error {
	v, ok := tunnels.Load(id)
	if !ok {
		return errors.New("通道不存在")
	}
	tunnels.Delete(id)
	return v.(Tunnel).Close()
}

func GetTunnel(id int) (Tunnel, error) {
	v, ok := tunnels.Load(id)
	if !ok {
		return nil, errors.New("通道不存在")
	}
	return v.(Tunnel), nil
}

func GetLink(tunnelId, linkId int) (Link, error) {
	t, err := GetTunnel(tunnelId)
	if err != nil {
		return nil, err
	}
	return t.GetLink(linkId)
}
