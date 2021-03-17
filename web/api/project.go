package api

import (
	"github.com/gin-gonic/gin"
)

func projectBeforeCreate(data interface{}) error {
	return nil
}

func projectAfterCreate(data interface{}) error {
	//project := data.(*model.Project)

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


func projectImport(c *gin.Context) {

}

func projectExport(c *gin.Context) {

}

func projectDeploy(c *gin.Context) {

}
