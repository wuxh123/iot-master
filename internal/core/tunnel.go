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
	GetTunnel() *models.ModelTunnel
	GetLink(id int64) (*Link, error)
}

type baseTunnel struct {
	models.Tunnel
	models.ModelTunnel
}

func (t *baseTunnel) GetTunnel() *models.ModelTunnel {
	return &t.ModelTunnel
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


func NewTunnel(tunnel *models.ModelTunnel) (Tunnel, error) {
	if tunnel.Role == "client" {
		return &TcpClient{
			baseTunnel: baseTunnel{
				ModelTunnel: *tunnel,
			},
		}, nil
	} else if tunnel.Role == "server" {
		switch tunnel.Net {
		case "tcp", "tcp4", "tcp6", "unix":
			return &TcpServer{
				baseTunnel: baseTunnel{
					ModelTunnel: *tunnel,
				},
			}, nil
		case "udp", "udp4", "udp6", "unixgram":
			return &PacketServer{
				baseTunnel: baseTunnel{
					ModelTunnel: *tunnel,
				},
			}, nil
		default:
			return nil, fmt.Errorf("未知的网络类型 %s", tunnel.Net)
		}
	} else {
		return nil, fmt.Errorf("未知的角色 %s", tunnel.Role)
	}
}
