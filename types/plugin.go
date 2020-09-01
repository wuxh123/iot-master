package types

import "time"

type Plugin struct {
	ID        int       `json:"id" storm:"increment"`
	Name      string    `json:"name"`
	AppKey    string    `json:"app_key" storm:"unique index"`
	AppSecret string    `json:"app_secret"`
	Address   string    `json:"address"`
	Path      string    `json:"path"`
	Entry     string    `json:"entry"`
	Created   time.Time `json:"created"`
}
