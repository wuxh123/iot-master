package storage

import "time"

type RegisterConf struct {
	Enable bool
	Length int
	Regex  string
}

type HeartBeatConf struct {
	Enable  bool
	Content []byte
}

type Channel struct {
	ID       int
	Serial   string `storm:"index"`
	Net      string
	Addr     string
	IsServer bool
	Name     string
	Tags     []string

	Register struct {
		Enable bool
		Length int
		Regex  string
	}

	HeartBeat struct {
		Enable   bool
		Interval time.Duration
		Content  []byte
	}

	Created time.Time
	Creator int
}

type User struct {
	ID       int
	Username string `storm:"index"`
	Password string
	Created  time.Time
}
