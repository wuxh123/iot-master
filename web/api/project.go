package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/gin-gonic/gin"
)

func projectImport(ctx *gin.Context) {
	var model models.Template

	err := ctx.ShouldBind(&model)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}

func projectExport(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var model models.Template
	has, err := db.Engine.ID(pid.Id).Table("model").Get(&model)
	if !has {
		replyFail(ctx, "记录不存在")
		return
	}
	if err != nil {
		replyError(ctx, err)
		return
	}

	//读取任务
	err = db.Engine.Where("model_id=?", pid.Id).Find(&model.Jobs)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//读取策略
	err = db.Engine.Where("model_id=?", pid.Id).Find(&model.Strategies)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, model)
}

func projectDeploy(ctx *gin.Context) {

}
