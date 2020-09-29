package api

import (
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/types"
	"github.com/gin-gonic/gin"
	"time"
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
	PollingInterval int    `json:"polling_interval"` //轮询间隔
}

type ModelVariable struct {
	ModelBase
	Tunnel   string `json:"tunnel"`
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
	Tunnel string `json:"tunnel"`
	Type   string `json:"type"`
	Addr   string `json:"addr"`
	Size   int    `json:"size"`

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`

	//结果解析
	Results []types.ModelBatchResult `json:"results"`
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
	H5          string `json:"h5"`
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
	m := types.Model{
		Name:        model.Name,
		Description: model.Description,
		Version:     model.Version,
		H5:          model.H5,
		Polling:     model.Polling,
		CreatedAt:   time.Now(),
	}

	modelDB := db.DB("model")
	err = modelDB.Save(&m)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//创建通道
	tunnelIds := make(map[string]int)
	tunnelDB := modelDB.From("tunnel")
	for _, t := range model.Tunnels {
		tunnel := types.ModelTunnel{
			ModelBase: types.ModelBase{
				ModelId:     m.Id,
				Name:        t.Name,
				Description: t.Description,
				CreatedAt:   time.Now(),
			},
			Protocol:        t.Protocol,
			ProtocolOpts:    t.ProtocolOpts,
			PollingEnable:   t.PollingEnable,
			PollingInterval: t.PollingInterval,
		}
		err = tunnelDB.Save(&tunnel)
		if err != nil {
			replyError(ctx, err)
			return
		}
		tunnelIds[tunnel.Name] = tunnel.Id
	}

	//创建变量
	variableDB := modelDB.From("variable")
	for _, v := range model.Variables {
		variable := types.ModelVariable{
			ModelBase: types.ModelBase{
				ModelId:     m.Id,
				Name:        v.Name,
				Description: v.Description,
				CreatedAt:   time.Now(),
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
		err = variableDB.Save(&variable)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}

	//创建批量
	batchDB := modelDB.From("batch")
	for _, v := range model.Batches {
		batch := types.ModelBatch{
			ModelBase: types.ModelBase{
				ModelId:     m.Id,
				Name:        v.Name,
				Description: v.Description,
				CreatedAt:   time.Now(),
			},
			TunnelId:      tunnelIds[v.Tunnel],
			Type:          v.Type,
			Addr:          v.Addr,
			Size:          v.Size,
			Cron:          v.Cron,
			PollingEnable: v.PollingEnable,
			PollingTimes:  v.PollingTimes,
			Results:       nil,
		}
		err = batchDB.Save(&batch)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}

	//创建任务
	jobDB := modelDB.From("job")
	for _, v := range model.Jobs {
		job := types.ModelJob{
			ModelBase: types.ModelBase{
				ModelId:     m.Id,
				Name:        v.Name,
				Description: v.Description,
				CreatedAt:   time.Now(),
			},
			Cron:   v.Cron,
			Script: v.Script,
		}
		err = jobDB.Save(&job)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}

	//创建策略
	strategyDB := modelDB.From("strategy")
	for _, v := range model.Strategies {
		job := types.ModelStrategy{
			ModelBase: types.ModelBase{
				ModelId:     m.Id,
				Name:        v.Name,
				Description: v.Description,
				CreatedAt:   time.Now(),
			},
			Script: v.Script,
		}
		err = strategyDB.Save(&job)
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

	var model types.Model
	modelDB := db.DB("model")
	err := modelDB.One("Id", pid.Id, &model)
	if err != nil {
		replyError(ctx, err)
		return
	}

	m := Model{
		Name:        model.Name,
		Description: model.Description,
		Version:     model.Version,
		H5:          model.H5,
		Polling:     model.Polling,
		Tunnels:     make([]ModelTunnel, 0),
		Variables:   make([]ModelVariable, 0),
		Batches:     make([]ModelBatch, 0),
		Jobs:        make([]ModelJob, 0),
		Strategies:  make([]ModelStrategy, 0),
	}

	//读取通道
	tunnelIds := make(map[int]string)
	tunnelDB := modelDB.From("tunnel")
	var tunnels []types.ModelTunnel
	err = tunnelDB.Find("ModelId", model.Id, &tunnels)
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
		}
		m.Tunnels = append(m.Tunnels, tunnel)
		tunnelIds[v.Id] = v.Name
	}

	//读取变量
	variableDB := modelDB.From("variable")
	var variables []types.ModelVariable
	err = variableDB.Find("ModelId", model.Id, &variables)
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
	batchDB := modelDB.From("batch")
	var batches []types.ModelBatch
	err = batchDB.Find("ModelId", model.Id, &batches)
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
			Results:       v.Results,
		}
		m.Batches = append(m.Batches, batch)
	}

	//读取任务
	jobDB := modelDB.From("job")
	var jobs []types.ModelJob
	err = jobDB.Find("ModelId", model.Id, &jobs)
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
	strategyDB := modelDB.From("strategy")
	var strategies []types.ModelStrategy
	err = strategyDB.Find("ModelId", model.Id, &strategies)
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
