package api

import (
	"git.zgwit.com/iot/dtu-admin/db"
	"git.zgwit.com/iot/dtu-admin/dtu"
	"git.zgwit.com/iot/dtu-admin/model"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/storm/v3/q"
	"log"
	"net/http"
)

func channels(ctx *gin.Context) {
	cs := make([]model.Channel, 0)

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
			q.Re("Addr", body.Keyword),
		))
	}

	query := db.DB("channel").Select(cond...)

	//计算总数
	cnt, err := query.Count(&model.Channel{})
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

func channelCreate(ctx *gin.Context) {
	var channel model.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("channel").Save(&channel)
	if err != nil {
		replyError(ctx, err)
		return
	}
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

	err := db.DB("channel").DeleteStruct(&model.Link{Id: pid.Id})
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
	err := db.DB("channel").Update(&channel)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, channel)

	//重新启动服务
	go func() {
		_ = dtu.DeleteChannel(channel.Id)
		//如果 disabled，则删除之
		if channel.Disabled {
			return
		}

		_, err := dtu.StartChannel(&channel)
		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func getChannelFromUri(ctx *gin.Context) (*model.Channel, error) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		return nil, err
	}

	var channel model.Channel
	err := db.DB("channel").One("Id", pid.Id, &channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func channelGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}
	var channel model.Channel
	err := db.DB("channel").One("Id", pid.Id, &channel)
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
