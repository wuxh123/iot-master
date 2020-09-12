package model

import "time"

type Plugin struct {
	Id        int       `json:"id" storm:"id,increment"`
	Name      string    `json:"name"`
	AppKey    string    `json:"app_key" store:"index"`
	AppSecret string    `json:"app_secret"`
	CreatedAt time.Time `json:"created_at"`
}
