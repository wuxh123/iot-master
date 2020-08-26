package dtu

import "sync"

type Link interface {
	Open() error
	Close() error
}

var links *sync.Map

func init() {
	links = new(sync.Map)
}

func Links() *sync.Map {
	return links
}