package core

import (
	"fmt"
	"git.zgwit.com/zgwit/iot-admin/models"
	"regexp"
)

//通道
type Tunnel interface {
	Open() error
	Close() error
	GetTunnel() *models.Tunnel
	GetLink(id int64) (*Link, error)
}

type baseTunnel struct {
	models.Tunnel
	//models.ProjectAdapter
}

func (t *baseTunnel) GetTunnel() *models.Tunnel {
	return &t.Tunnel
}

func (t *baseTunnel) checkRegister(buf []byte) (string, error) {
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

func NewTunnel(tunnel *models.Tunnel) (Tunnel, error) {
	switch tunnel.Type {
	case "tcp-server":
		return &TcpServer{
			baseTunnel: baseTunnel{
				Tunnel: *tunnel,
			},
		}, nil
	case "tcp-client":
		return &TcpUdpClient{
			baseTunnel: baseTunnel{
				Tunnel: *tunnel,
			},
			Net: "tcp",
		}, nil
	case "udp-server":
		return &PacketServer{
			baseTunnel: baseTunnel{
				Tunnel: *tunnel,
			},
		}, nil
	case "udp-client":
		return &TcpUdpClient{
			baseTunnel: baseTunnel{
				Tunnel: *tunnel,
			},
			Net: "udp",
		}, nil
	case "serial":
	default:
	}
	return nil, fmt.Errorf("未知的网络类型 %s", tunnel.Type)
}
