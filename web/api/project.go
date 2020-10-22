package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/kataras/iris/v12"
)

func projectImport(ctx iris.Context) {
	var model models.Template

	err := ctx.ReadJSON(&model)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}

func projectExport(ctx iris.Context) {
	id, err := ctx.URLParamInt64("id")
	if err != nil {
		replyError(ctx, err)
		return
	}

	var model models.Template
	has, err := db.Engine.ID(id).Table("model").Get(&model)
	if !has {
		replyFail(ctx, "记录不存在")
		return
	}
	if err != nil {
		replyError(ctx, err)
		return
	}

	//读取任务
	err = db.Engine.Where("model_id=?", id).Find(&model.Jobs)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//读取策略
	err = db.Engine.Where("model_id=?", id).Find(&model.Strategies)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, model)
}

func projectDeploy(ctx iris.Context) {

}
