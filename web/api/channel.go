package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/model"
	"log"
	"net/http"
)

func channels(ctx *gin.Context) {
	cs := make([]model.Channel, 0)

	var body paramSearch
	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//op := db.Engine.Where("type=?", body.Net)
	op := db.Engine.NewSession()
	for _, filter := range body.Filters {
		if len(filter.Value) > 0 {
			if len(filter.Value) == 1 {
				op.And(filter.Key+"=?", filter.Value[0])
			} else {
				op.In(filter.Key, filter.Value)
			}
		}
	}
	if body.Keyword != "" {
		kw := "%" + body.Keyword + "%"
		op.And("name like ? or addr like ?", kw, kw)
	}

	op.Limit(body.Length, body.Offset)
	if body.SortKey != "" {
		if body.SortOrder == "desc" {
			op.Desc(body.SortKey)
		} else {
			op.Asc(body.SortKey)
		}
	} else {
		op.Desc("id")
	}
	cnt, err := op.FindAndCount(&cs)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//replyOk(ctx, cs)
	ctx.JSON(http.StatusOK, gin.H{
		"ok":    true,
		"data":  cs,
		"total": cnt,
	})
}

func channelCreate(ctx *gin.Context) {
	var channel model.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		replyError(ctx, err)
		return
	}

	// channel.Creator = TODO 从session中获取

	_, err := db.Engine.Insert(&channel)
	if err != nil {
		replyError(ctx, err)
		return
	}
	//获取完整内容
	_, _ = db.Engine.ID(channel.Id).Get(&channel)
	replyOk(ctx, channel)

	//启动服务
	go func() {
		_, err := dtu.StartChannel(&channel)
		if err != nil {
			log.Println(err)
		}
	}()
}

func channelDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	_, err := db.Engine.ID(pid.Id).Get(&model.Channel{})
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, nil)

	//删除服务
	go func() {
		channel, err := dtu.GetChannel(pid.Id)
		if err != nil {
			log.Println(err)
			return
		}

		err = channel.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func channelModify(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var channel model.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		replyError(ctx, err)
		return
	}

	//log.Println("update", channel)
	_, err := db.Engine.ID(pid.Id).
		Cols("name", "disabled",
			"type", "addr", "role", "timeout",
			"register_enable", "register_regex",
			"heart_beat_enable", "heart_beat_interval", "heart_beat_content", "heart_beat_is_hex",
			"plugin_id").Update(&channel)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, channel)

	//重新启动服务
	go func() {
		//如果 disabled，则删除之
		if channel.Disabled {
			_ = dtu.DeleteChannel(channel.Id)
			return
		}

		ch, err := dtu.GetChannel(pid.Id)
		if err != nil {
			log.Println(err)
			//找不到，重新启动一个
			_, _ = dtu.StartChannel(&channel)
			return
		}

		err = ch.Close()
		if err != nil {
			log.Println(err)
			return
		}

		//重置参数，重新启动
		ch.Channel = channel
		err = ch.Open()
		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func channelGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var channel model.Channel
	has, err := db.Engine.ID(pid.Id).Get(&channel)
	if err != nil {
		replyError(ctx, err)
		return
	}
	if !has {
		replyFail(ctx, "找不到通道")
		return
	}

	replyOk(ctx, channel)
}

func channelStart(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}
	c, err := dtu.GetChannel(pid.Id)
	if err != nil {
		replyError(ctx, err)
		return
	}

	err = c.Open()
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}

func channelStop(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}
	c, err := dtu.GetChannel(pid.Id)
	if err != nil {
		replyError(ctx, err)
		return
	}

	err = c.Close()
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}
