package api

import (
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/google/uuid"
	"net/http"
)

func projectBeforeCreate(data interface{}) error {
	project := data.(*models.Project)
	project.UUID = uuid.New().String()
	return nil
}

func projectAfterCreate(data interface{}) error {
	//project := data.(*models.Project)

	//TODO 加载实例
	return nil
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
