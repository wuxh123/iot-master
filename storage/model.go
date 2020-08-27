package storage

import "time"

type Channel struct {
	ID       int
	Serial   string `storm:"index"`
	Net      string
	Addr     string
	IsServer bool
	Name     string
	Tags     []string
	Created  time.Time
	Creator  int
}

type User struct {
	ID       int
	Username string `storm:"index"`
	Password string
	Created  time.Time
}
