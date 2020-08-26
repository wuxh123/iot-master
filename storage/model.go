package storage

import "time"

type Device struct {
	ID        int
	Serial    string `storm:"index"`
	Name      string
	Tags      []string
	CreatedAt time.Time
}