package protocol

import (
	"errors"
	"sync"
)

//可以改为普通map
var protocols sync.Map //protocol

//功能码
type Code struct {
	Name string `json:"name"`
	Code uint8  `json:"code"`
}

type protocol struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Codes   []Code `json:"codes"`
	factory Factory
}

type Factory func(opts string) (Adapter, error)

func GetProtocols() []protocol {
	ps := make([]protocol, 0)
	protocols.Range(func(key, value interface{}) bool {
		ps = append(ps, value.(protocol))
		return true
	})
	return ps
}

func CreateAdapter(name string, opts string) (Adapter, error) {
	if v, ok := protocols.Load(name); ok && v != nil {
		p := v.(*protocol)
		return p.factory(opts)
	}
	return nil, errors.New("找不到协议")
}

func RegisterAdapter(name, version string, codes []Code, factory Factory) {
	protocols.Store(name, &protocol{
		Name:    name,
		Codes:   codes,
		factory: factory,
	})
}
