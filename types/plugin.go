package types

import "time"

type Plugin struct {
	ID        int       `storm:"increment" json:"id"`
	Name      string    `json:"name"`
	AppKey    string    `storm:"unique index" json:"app_key"`
	AppSecret string    `json:"app_secret"`
	Address   string    `json:"address"`
	Path      string    `json:"path"`
	Entry     string    `json:"entry"`
	Created   time.Time `json:"created"`
}
