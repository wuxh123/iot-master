package storage

import (
	"github.com/asdine/storm/v3"
)

var (
	DeviceDB *storm.DB
)

func Open() error {
	var err error
	DeviceDB, err = storm.Open("device.db")
	if err != nil {
		return err
	}

	return nil
}