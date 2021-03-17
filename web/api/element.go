package api

import (
	"github.com/google/uuid"
	"iot-master/model"
)

func elementBeforeCreate(data interface{}) error {
	element := data.(*model.Element)
	if element.Origin == "" {
		element.Origin = uuid.New().String()
	}
	return nil
}

func elementBeforeDelete(data interface{}) error {
	//TODO 检查是否被引用

	return nil
}
