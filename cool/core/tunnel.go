package core

import (
	"fmt"
	"git.zgwit.com/zgwit/iot-admin/models"
	"regexp"
)

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

