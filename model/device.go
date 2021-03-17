package model

import "time"

type Device struct {
	Id        int64 `json:"id"`
	ProjectId int64   `json:"project_id"`

	Created time.Time `json:"created" storm:"created"`
}
