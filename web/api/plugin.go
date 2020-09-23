package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"git.zgwit.com/zgwit/iot-admin/model"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/storm/v3/q"
	"net/http"
)

func plugins(ctx *gin.Context) {
	cs := make([]model.Plugin, 0)

	var body paramSearch
	err := ctx.ShouldBind(&body)
	if err != nil {
		replyError(ctx, err)
		return
	}

	cond := make([]q.Matcher, 0)
	//过滤条件
	for _, filter := range body.Filters {
		if len(filter.Value) > 0 {
			cond = append(cond, q.In(filter.Key, filter.Value))
		}
	}
	//关键字
	if body.Keyword != "" {
		cond = append(cond, q.Or(
			q.Re("Name", body.Keyword),
			q.Re("Key", body.Keyword),
		))
	}

	query := db.DB("plugin").Select(cond...)

	//计算总数
	cnt, err := query.Count(&model.Plugin{})
	if err != nil {
		replyError(ctx, err)
		return
	}

	//分页
	query = query.Skip(body.Offset).Limit(body.Length)

	//排序
	if body.SortKey != "" {
		if body.SortOrder == "desc" {
			query = query.OrderBy(body.SortKey).Reverse()
		} else {
			query = query.OrderBy(body.SortKey)
		}
	} else {
		query = query.OrderBy("Id").Reverse()
	}

	err = query.Find(&cs)
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

func pluginCreate(ctx *gin.Context) {
	var plugin model.Plugin
	if err := ctx.ShouldBindJSON(&plugin); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("plugin").Save(&plugin)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, plugin)
}

func pluginDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("plugin").DeleteStruct(&model.Link{Id: pid.Id})
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, nil)
}

func pluginModify(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var plugin model.Plugin
	if err := ctx.ShouldBindJSON(&plugin); err != nil {
		replyError(ctx, err)
		return
	}

	//log.Println("update", plugin)
	err := db.DB("plugin").Update(&plugin)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, plugin)
}


func pluginGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}
	var plugin model.Plugin
	err := db.DB("plugin").One("Id", pid.Id, &plugin)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, plugin)
}
