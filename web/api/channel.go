package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/storage"
	"github.com/zgwit/dtu-admin/types"
)

func channelAll(ctx *gin.Context) {
	var cs []types.Channel
	err := storage.DB("channel").All(&cs)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, cs)
}

func channels(ctx *gin.Context) {
	cs := dtu.Channels()
	replyOk(ctx, cs)
}

func channelCreate(ctx *gin.Context) {
	var channel types.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		replyError(ctx, err)
		return
	}

	//创建并启动
	c, err := dtu.CreateChannel(&channel)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, c)
}

func channelDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}
	err := dtu.DeleteChannel(pid.Id)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}

func channelModify(ctx *gin.Context) {

}

func channelGet(ctx *gin.Context) {
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

	replyOk(ctx, c)
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
