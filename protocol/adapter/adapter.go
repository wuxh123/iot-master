package adapter

import (
	"errors"
	"git.zgwit.com/iot/mydtu/base"
	"sync"
)

type Adapter interface {
	Name() string
	Version() string

	Attach(link base.Link)

	Read(slave uint8, code uint8, offset uint16, size uint16) ([]byte, error)
	Write(slave uint8, code uint8, offset uint16, buf []byte) error
}

type Factory func(opts string) (Adapter, error)

type Code struct {
	Name string `json:"name"`
	Code uint8  `json:"code"`
}

type adapter struct {
	Name    string `json:"name"`
	Codes   []Code `json:"codes"`
	factory Factory
}

//可以改为普通map
var protocols sync.Map //adapter

func GetProtocols() []adapter {
	ps := make([]adapter, 0)
	protocols.Range(func(key, value interface{}) bool {
		ps = append(ps, value.(adapter))
		return true
	})
	return ps
}

func CreateAdapter(name string, opts string) (Adapter, error) {
	if v, ok := protocols.Load(name); ok && v != nil {
		p := v.(*adapter)
		return p.factory(opts)
	}
	return nil, errors.New("找不到协议")
}

func RegisterAdapter(name string, codes []Code, factory Factory) {
	protocols.Store(name, &adapter{
		Name:    name,
		Codes:   codes,
		factory: factory,
	})
}
