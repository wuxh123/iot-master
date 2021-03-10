package model

import "time"

type Plugin struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Key      string    `json:"key" storm:"unique"`
	Secret   string    `json:"secret"`
	Disabled bool      `json:"disabled"`
	Created  time.Time `json:"created" storm:"created"`
	Updated  time.Time `json:"updated" storm:"updated"`
}
