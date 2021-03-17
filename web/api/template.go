package api

import (
	"github.com/google/uuid"
	"iot-master/model"
)

func templateBeforeCreate(data interface{}) error {
	template := data.(*model.Template)
	if template.Origin == "" {
		template.Origin = uuid.New().String()
	}
	return nil
}

