package api

import (
	"iot-master/model"
	"github.com/google/uuid"
)

func elementBeforeCreate(data interface{}) error {
	element := data.(*model.Element)
	element.UUID = uuid.New().String()
	return nil
}

func elementBeforeDelete(data interface{}) error {
	//TODO 检查是否被引用

	return nil
}
