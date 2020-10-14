package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/gin-gonic/gin"
)

type ModelBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModelTunnel struct {
	ModelBase
	Protocol        string `json:"protocol"`
	ProtocolOpts    string `json:"protocol_opts"`
	PollingEnable   bool   `json:"polling_enable"`   //轮询
	PollingInterval int    `json:"polling_interval"` //轮询间隔 ms
	PollingCycle    int    `json:"polling_cycle"`    //轮询周期 s
}

type ModelVariable struct {
	ModelBase
	Tunnel   string `json:"core"`
	Type     string `json:"type"`
	Addr     string `json:"addr"`
	Default  string `json:"default"`
	Writable bool   `json:"writable"` //可写，用于输出（如开关）

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`
}

type ModelBatch struct {
	ModelBase
	Tunnel string `json:"core"`
	Type   string `json:"type"`
	Addr   string `json:"addr"`
	Size   int    `json:"size"`

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`

	//结果解析
	Results []models.ModelBatchResult `json:"results"`
}

type ModelJob struct {
	ModelBase
	Cron   string `json:"cron"`
	Script string `json:"script"` //javascript
}

type ModelStrategy struct {
	ModelBase
	Script string `json:"script"` //javascript
}

type Model struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Polling     bool   `json:"polling"` //轮询

	Tunnels    []ModelTunnel   `json:"tunnels"`
	Variables  []ModelVariable `json:"variables"`
	Batches    []ModelBatch    `json:"batches"`
	Jobs       []ModelJob      `json:"jobs"`
	Strategies []ModelStrategy `json:"strategies"`
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
	tunnelIds := make(map[string]int64)
	for _, t := range model.Tunnels {
		tunnel := models.ModelTunnel{
			ModelBase: models.ModelBase{
				ModelId:     m.Id,
				Name:        t.Name,
				Description: t.Description,
			},
			Protocol:        t.Protocol,
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
		tunnelIds[tunnel.Name] = tunnel.Id
	}

	//创建变量
	for _, v := range model.Variables {
		variable := models.ModelVariable{
			ModelBase: models.ModelBase{
				ModelId:     m.Id,
				Name:        v.Name,
				Description: v.Description,
			},
			TunnelId:      tunnelIds[v.Tunnel],
			Type:          v.Type,
			Addr:          v.Addr,
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
			TunnelId:      tunnelIds[v.Tunnel],
			Type:          v.Type,
			Addr:          v.Addr,
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

	var model models.Model
	has, err := db.Engine.ID(pid.Id).Get(&model)
	if !has {
		replyFail(ctx, "记录不存在")
		return
	}
	if err != nil {
		replyError(ctx, err)
		return
	}

	m := Model{
		Name:        model.Name,
		Description: model.Description,
		Version:     model.Version,
		Tunnels:     make([]ModelTunnel, 0),
		Variables:   make([]ModelVariable, 0),
		Batches:     make([]ModelBatch, 0),
		Jobs:        make([]ModelJob, 0),
		Strategies:  make([]ModelStrategy, 0),
	}

	//读取通道
	tunnelIds := make(map[int64]string)
	var tunnels []models.ModelTunnel
	err = db.Engine.Where("model_id=?", model.Id).Find(&tunnels)
	if err != nil {
		replyError(ctx, err)
		return
	}
	for _, v := range tunnels {
		tunnel := ModelTunnel{
			ModelBase: ModelBase{
				Name:        v.Name,
				Description: v.Description,
			},
			Protocol:        v.Protocol,
			ProtocolOpts:    v.ProtocolOpts,
			PollingEnable:   v.PollingEnable,
			PollingInterval: v.PollingInterval,
			PollingCycle:    v.PollingCycle,
		}
		m.Tunnels = append(m.Tunnels, tunnel)
		tunnelIds[v.Id] = v.Name
	}

	//读取变量
	var variables []models.ModelVariable
	err = db.Engine.Where("model_id=?", model.Id).Find(&variables)
	if err != nil {
		replyError(ctx, err)
		return
	}
	for _, v := range variables {
		variable := ModelVariable{
			ModelBase: ModelBase{
				Name:        v.Name,
				Description: v.Description,
			},
			Tunnel:        tunnelIds[v.TunnelId],
			Type:          v.Type,
			Addr:          v.Addr,
			Default:       v.Default,
			Writable:      v.Writable,
			Cron:          v.Cron,
			PollingEnable: v.PollingEnable,
			PollingTimes:  v.PollingTimes,
		}
		m.Variables = append(m.Variables, variable)
	}

	//读取批量
	var batches []models.ModelBatch
	err = db.Engine.Where("model_id=?", model.Id).Find(&batches)
	if err != nil {
		replyError(ctx, err)
		return
	}
	for _, v := range batches {
		batch := ModelBatch{
			ModelBase: ModelBase{
				Name:        v.Name,
				Description: v.Description,
			},
			Tunnel:        tunnelIds[v.TunnelId],
			Type:          v.Type,
			Addr:          v.Addr,
			Size:          v.Size,
			Cron:          v.Cron,
			PollingEnable: v.PollingEnable,
			PollingTimes:  v.PollingTimes,
			//Results:       v.Results, //TODO results
		}
		m.Batches = append(m.Batches, batch)
	}

	//读取任务
	var jobs []models.ModelJob
	err = db.Engine.Where("model_id=?", model.Id).Find(&jobs)
	if err != nil {
		replyError(ctx, err)
		return
	}
	for _, v := range jobs {
		job := ModelJob{
			ModelBase: ModelBase{
				Name:        v.Name,
				Description: v.Description,
			},
			Cron:   v.Cron,
			Script: v.Script,
		}
		m.Jobs = append(m.Jobs, job)
	}

	//读取策略
	var strategies []models.ModelStrategy
	err = db.Engine.Where("model_id=?", model.Id).Find(&strategies)
	if err != nil {
		replyError(ctx, err)
		return
	}
	for _, v := range strategies {
		strategy := ModelStrategy{
			ModelBase: ModelBase{
				Name:        v.Name,
				Description: v.Description,
			},
			Script: v.Script,
		}
		m.Strategies = append(m.Strategies, strategy)
	}

	replyOk(ctx, m)
}

func modelRefresh(ctx *gin.Context) {

}