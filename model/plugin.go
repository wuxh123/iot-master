package model

import "time"

type Plugin struct {
	Id        int       `json:"id" storm:"id,increment"`
	Name      string    `json:"name"`
	Key       string    `json:"key" store:"unique"`
	Secret    string    `json:"secret"`
	CreatedAt time.Time `json:"created_at"`
}
