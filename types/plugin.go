package types

import "time"

type Plugin struct {
	Id        int       `json:"id" storm:"id,increment"`
	Name      string    `json:"name"`
	Key       string    `json:"key" storm:"unique"`
	Secret    string    `json:"secret"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
}
