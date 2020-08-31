package types

import "time"

type User struct {
	ID       int       `storm:"increment" json:"id"`
	Username string    `storm:"unique index" json:"username"`
	Password string    `json:"password"`
	Disabled bool      `json:"disabled"` //TODO 未实现
	Created  time.Time `json:"created"`
}
