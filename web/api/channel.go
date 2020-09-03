package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/storage"
	"github.com/zgwit/dtu-admin/model"
	"log"
	"time"
)


func channels(ctx *gin.Context) {
	var cs []model.Channel
	err := storage.DB("channel").All(&cs)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, cs)
}

func channelCreate(ctx *gin.Context) {
	var channel model.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		replyError(ctx, err)
		return
	}

	// channel.Creator = TODO 从session中获取
	channel.Created = time.Now()
	err := storage.DB("channel").Save(&channel)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//启动服务
	go func() {
		_, err := dtu.StartChannel(&channel)
		if err != nil {
			log.Println(err)
		}
	}()

	replyOk(ctx, channel)
}

func channelDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	err := storage.DB("channel").DeleteStruct(&model.Channel{ID: pid.Id})
	if err != nil {
		replyError(ctx, err)
		return
	}

	//TODO 删除服务


	replyOk(ctx, nil)
}


func channelModify(ctx *gin.Context) {
	var channel model.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		replyError(ctx, err)
		return
	}

	log.Println("update", channel)

	//TODO 不能全部字段更新，应该先取值，修改，再存入
	err := storage.DB("channel").Update(&channel)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//TODO 重新启动服务

	replyOk(ctx, channel)
}

func channelGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var channel model.Channel
	err := storage.DB("channel").One("ID", pid.Id, &channel)
	if err != nil {
		replyError(ctx, err)
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
