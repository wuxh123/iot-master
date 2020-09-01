package types

import "time"

type User struct {
	ID       int       `json:"id" storm:"increment"`
	Username string    `json:"username" storm:"unique"`
	Password string    `json:"password"`
	Disabled bool      `json:"disabled"` //TODO 未实现
	Created  time.Time `json:"created"`
}
