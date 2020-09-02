package types

import (
	"time"
)

type Link struct {
	ID      int       `json:"id" storm:"increment"`
	Name    string    `json:"name"`
	Serial  string    `json:"serial" storm:"index"`
	Addr    string    `json:"addr"`
	Channel int       `json:"channel"`
	Online  time.Time `json:"online"`
	Created time.Time `json:"created"`
}
