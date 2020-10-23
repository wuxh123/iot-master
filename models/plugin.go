package models

import "time"

type Plugin struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Key      string    `json:"key" xorm:"unique"`
	Secret   string    `json:"secret"`
	Disabled bool      `json:"disabled"`
	Created  time.Time `json:"created" xorm:"created"`
	Updated  time.Time `json:"updated" xorm:"updated"`
}
