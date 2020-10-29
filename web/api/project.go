package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"github.com/google/uuid"
	"net/http"
)

func projectAfterCreate(data interface{}) error {
	//project := data.(*models.Project)
	return db.DB("project").UpdateField(data, "UUID", uuid.New().String())

	//TODO 加载实例
}

func projectAfterModify(data interface{}) error {
	//TODO 修改实例
	return nil
}


func projectAfterDelete(data interface{}) error {
	//TODO 删除实例
	return nil
}


func projectImport(writer http.ResponseWriter, request *http.Request) {

}

func projectExport(writer http.ResponseWriter, request *http.Request) {

}

func projectDeploy(writer http.ResponseWriter, request *http.Request) {

}
