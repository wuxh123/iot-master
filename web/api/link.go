package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/dtu-admin/storage"
	"github.com/zgwit/dtu-admin/model"
	"log"
)

func links(ctx *gin.Context) {
	var cs []model.Link
	err := storage.DB("link").All(&cs)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, cs)
}

func linkDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	err := storage.DB("link").DeleteStruct(&model.Link{ID: pid.Id})
	if err != nil {
		replyError(ctx, err)
		return
	}

	//TODO 删除服务


	replyOk(ctx, nil)
}


func linkModify(ctx *gin.Context) {
	var link model.Link
	if err := ctx.ShouldBindJSON(&link); err != nil {
		replyError(ctx, err)
		return
	}

	log.Println("update", link)

	//TODO 不能全部字段更新，应该先取值，修改，再存入
	err := storage.DB("link").Update(&link)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//TODO 重新启动服务

	replyOk(ctx, link)
}

func linkGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	err := storage.DB("link").One("ID", pid.Id, &link)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, link)
}
