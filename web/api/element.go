package api

import (
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/google/uuid"
)

func elementBeforeCreate(data interface{}) error {
	element := data.(*models.Element)
	element.UUID = uuid.New().String()
	return nil
}

func elementBeforeDelete(data interface{}) error {
	//TODO 检查是否被引用

	return nil
}
