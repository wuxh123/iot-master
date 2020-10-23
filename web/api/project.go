package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"github.com/google/uuid"
	"net/http"
)

func projectAfterCreate(data interface{}) error {
	//project := data.(*models.Project)
	return db.DB("project").UpdateField(data, "UUID", uuid.New().String())
}

func projectImport(writer http.ResponseWriter, request *http.Request) {

}

func projectExport(writer http.ResponseWriter, request *http.Request) {

}

func projectDeploy(writer http.ResponseWriter, request *http.Request) {

}
