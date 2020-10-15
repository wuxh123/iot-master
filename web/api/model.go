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
	Polling     bool   `json:"polling"` //轮询

	Tunnels    []ModelTunnel   `json:"tunnels"`
	Variables  []ModelVariable `json:"variables"`
	Batches    []ModelBatch    `json:"batches"`
	Jobs       []ModelJob      `json:"jobs"`
	Strategies []ModelStrategy `json:"strategies"`
}

type ModelBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModelTunnel struct {
	ModelBase

	Role    string `json:"role"`
	Net     string `json:"net"`
	Addr    string `json:"addr"`
	Timeout int    `json:"timeout"`

	RegisterEnable bool   `json:"register_enable"`
	RegisterRegex  string `json:"register_regex"`
	RegisterMin    int    `json:"register_min"`
	RegisterMax    int    `json:"register_max"`

	HeartBeatEnable   bool   `json:"heart_beat_enable"`
	HeartBeatInterval int    `json:"heart_beat_interval"`
	HeartBeatContent  string `json:"heart_beat_content"`
	HeartBeatIsHex    bool   `json:"heart_beat_is_hex"`

	ProtocolName string `json:"protocol"`
	ProtocolOpts string `json:"protocol_opts"`

	PollingEnable   bool `json:"polling_enable"`   //轮询
	PollingInterval int  `json:"polling_interval"` //轮询间隔 ms
	PollingCycle    int  `json:"polling_cycle"`    //轮询周期 s
}

type ModelVariable struct {
	ModelBase

	Tunnel string `json:"tunnel"`
	models.Address

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
	ModelBase

	Tunnel string `json:"tunnel"`

	models.Address `xorm:"extends"`

	Size int `json:"size"`

	//采样：无、定时、轮询
	Cron          string `json:"cron"`
	PollingEnable bool   `json:"polling_enable"` //轮询
	PollingTimes  int    `json:"polling_times"`
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
			TunnelId:      tunnelIds[v.Tunnel],
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
	var tunnels []models.ModelAdapter
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
			ProtocolName:    v.ProtocolName,
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
			Address:       v.Address,
			Unit:          v.Unit,
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
			Address:       v.Address,
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
