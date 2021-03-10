package model

import (
	"time"
)

type User struct {
	Id       int64     `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	Disabled bool      `json:"disabled"`
	Created  time.Time `json:"created" storm:"created"`
}
