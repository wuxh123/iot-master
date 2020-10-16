package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/gin-gonic/gin"
)

type Model struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Adapter    ModelAdapter    `json:"adapter"`
	Variables  []ModelVariable `json:"variables"`
	Batches    []ModelBatch    `json:"batches"`
	Jobs       []ModelJob      `json:"jobs"`
	Strategies []ModelStrategy `json:"strategies"`
}

type ModelBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModelAdapter struct {
	ModelBase `xorm:"extends"`

	ProtocolName string `json:"protocol_name"`
	ProtocolOpts string `json:"protocol_opts"`

	PollingEnable   bool `json:"polling_enable"`   //轮询
	PollingInterval int  `json:"polling_interval"` //轮询间隔 ms
	PollingCycle    int  `json:"polling_cycle"`    //轮询周期 s
}

type ModelVariable struct {
	ModelBase `xorm:"extends"`

	models.Address `xorm:"extends"`

	Type string `json:"type"`
	Unit string `json:"unit"` //单位

	Scale    float32 `json:"scale"` //倍率，比如一般是 整数÷10，得到
	Default  string  `json:"default"`
	Writable bool    `json:"writable"` //可写，用于输出（如开关）

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`
}

type ModelBatch struct {
	ModelBase `xorm:"extends"`

	models.Address `xorm:"extends"`

	Size int `json:"size"`

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`
}

type ModelJob struct {
	ModelBase `xorm:"extends"`

	Cron   string `json:"cron"`
	Script string `json:"script"` //javascript
}

type ModelStrategy struct {
	ModelBase `xorm:"extends"`

	Script string `json:"script"` //javascript
}

func modelImport(ctx *gin.Context) {
	var model Model

	err := ctx.ShouldBind(&model)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//TODO 导入模型
	m := models.Model{
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
	tunnel := models.ModelAdapter{
		ModelBase: models.ModelBase{
			ModelId:     m.Id,
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
		variable := models.ModelVariable{
			ModelBase: models.ModelBase{
				ModelId:     m.Id,
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
		batch := models.ModelBatch{
			ModelBase: models.ModelBase{
				ModelId:     m.Id,
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
		job := models.ModelJob{
			ModelBase: models.ModelBase{
				ModelId:     m.Id,
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
		strategy := models.ModelStrategy{
			ModelBase: models.ModelBase{
				ModelId:     m.Id,
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

func modelExport(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var model Model
	has, err := db.Engine.ID(pid.Id).Table("model").Get(&model)
	if !has {
		replyFail(ctx, "记录不存在")
		return
	}
	if err != nil {
		replyError(ctx, err)
		return
	}

	//model := Model{
	//	Variables:   make([]ModelVariable, 0),
	//	Batches:     make([]ModelBatch, 0),
	//	Jobs:        make([]ModelJob, 0),
	//	Strategies:  make([]ModelStrategy, 0),
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

func modelRefresh(ctx *gin.Context) {

}
