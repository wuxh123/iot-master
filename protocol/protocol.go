package protocol

import (
	"errors"
	"log"
	"sync"
)

//可以改为普通map
var protocols sync.Map //Protocol

//功能码
type Code struct {
	Name string `json:"name"`
	Code uint8  `json:"code"`
}

//协议说明
type Manifest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Codes   []Code `json:"codes"`
}

type Protocol struct {
	Manifest *Manifest `json:"manifest"`
	factory  Factory
}

type Factory func(opts string) (Adapter, error)

func GetProtocols() []Protocol {
	ps := make([]Protocol, 0)
	protocols.Range(func(key, value interface{}) bool {
		ps = append(ps, value.(Protocol))
		return true
	})
	return ps
}

func CreateAdapter(name string, opts string) (Adapter, error) {
	if v, ok := protocols.Load(name); ok && v != nil {
		p := v.(*Protocol)
		return p.factory(opts)
	}
	return nil, errors.New("找不到协议")
}

func RegisterProtocol(name string, manifest *Manifest, factory Factory) {
	log.Println("加载协议", name)
	protocols.Store(name, &Protocol{
		Manifest: manifest,
		factory:  factory,
	})
}
