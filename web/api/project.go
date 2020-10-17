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

	//TODO 导入模型
	m := models.Project{
		Name:        model.Name,
		Description: model.Description,
		Version:     model.Version,
	}

	//TODO 根据origin查重
	_, err = db.Engine.Insert(&m)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//创建通道	
	t := model.Adapter
	tunnel := models.ProjectAdapter{
		ProjectBase: models.ProjectBase{
			ProjectId:   m.Id,
			Name:        t.Name,
			Description: t.Description,
		},
		ProtocolName:    t.ProtocolName,
		ProtocolOpts:    t.ProtocolOpts,
		PollingEnable:   t.PollingEnable,
		PollingInterval: t.PollingInterval,
		PollingCycle:    t.PollingCycle,
	}
	_, err = db.Engine.Insert(&tunnel)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//创建变量
	for _, v := range model.Variables {
		variable := models.ProjectVariable{
			ProjectBase: models.ProjectBase{
				ProjectId:   m.Id,
				Name:        v.Name,
				Description: v.Description,
			},
			Type:          v.Type,
			Address:       v.Address,
			Default:       v.Default,
			Writable:      v.Writable,
			Cron:          v.Cron,
			PollingEnable: v.PollingEnable,
			PollingTimes:  v.PollingTimes,
		}
		_, err = db.Engine.Insert(&variable)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}

	//创建批量
	for _, v := range model.Batches {
		batch := models.ProjectBatch{
			ProjectBase: models.ProjectBase{
				ProjectId:   m.Id,
				Name:        v.Name,
				Description: v.Description,
			},
			Address:       v.Address,
			Size:          v.Size,
			Cron:          v.Cron,
			PollingEnable: v.PollingEnable,
			PollingTimes:  v.PollingTimes,
		}
		_, err = db.Engine.Insert(&batch)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}

	//创建任务
	for _, v := range model.Jobs {
		job := models.ProjectJob{
			ProjectBase: models.ProjectBase{
				ProjectId:   m.Id,
				Name:        v.Name,
				Description: v.Description,
			},
			Cron:   v.Cron,
			Script: v.Script,
		}
		_, err = db.Engine.Insert(&job)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}

	//创建策略
	for _, v := range model.Strategies {
		strategy := models.ProjectStrategy{
			ProjectBase: models.ProjectBase{
				ProjectId:   m.Id,
				Name:        v.Name,
				Description: v.Description,
			},
			Script: v.Script,
		}
		_, err = db.Engine.Insert(&strategy)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}

	replyOk(ctx, m)
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

	//model := Template{
	//	Variables:   make([]ProjectVariable, 0),
	//	Batches:     make([]ProjectBatch, 0),
	//	Jobs:        make([]ProjectJob, 0),
	//	Strategies:  make([]ProjectStrategy, 0),
	//}

	//读取通道
	has, err = db.Engine.Where("model_id=?", pid.Id).Table("model_adapter").Get(&model.Adapter)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//读取变量
	err = db.Engine.Where("model_id=?", pid.Id).Find(&model.Variables)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//读取批量
	err = db.Engine.Where("model_id=?", pid.Id).Find(&model.Batches)
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
