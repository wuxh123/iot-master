package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"github.com/google/uuid"
)

func elementAfterCreated(data interface{}) error {
	//element := data.(*models.Element)
	return db.DB("element").UpdateField(data, "UUID", uuid.New().String())
}
