package dtu

import (
	"github.com/zgwit/dtu-admin/storage"
	"sync"
)

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

func Recovery() error  {
	var links []storage.Link
	err := storage.Links.All(&links)
	if err != nil {
		return err
	}

	for _, v := range links {
		if v.IsServer {
			lnk := NewServer(v.Net, v.Addr)
			err := lnk.Open()
			if err != nil {
				lnk.Err = err.Error()
			}
		} else {
			lnk := NewClient(v.Net, v.Addr)
			err := lnk.Open()
			if err != nil {
				lnk.Err = err.Error()
			}
		}
	}

	return nil
}