package models

import "time"

type Device struct {
	ID        int `json:"id"`
	ProjectId int `json:"project_id"`

	Created time.Time `json:"created" storm:"created"`
}
