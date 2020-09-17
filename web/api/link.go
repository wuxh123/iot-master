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

func links(ctx *gin.Context) {
	var ls []model.Link

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
			q.Re("Serial", body.Keyword),
			q.Re("Addr", body.Keyword),
		))
	}

	query := db.DB("link").Select(cond...)

	//计算总数
	cnt, err := query.Count(&model.Link{})
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

	err = query.Find(&ls)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//replyOk(ctx, cs)
	ctx.JSON(http.StatusOK, gin.H{
		"ok":    true,
		"data":  ls,
		"total": cnt,
	})
}

func linkDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	err := db.DB("link").DeleteStruct(&model.Link{Id: pid.Id})
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, nil)

	//删除服务
	go func() {
		l, err := dtu.GetLink(link.ChannelId, link.Id)
		if err != nil {
			log.Println(err)
			return
		}
		_ = l.Close()
		//TODO 强制删除连接
	}()

}

func linkModify(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	if err := ctx.ShouldBindJSON(&link); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("link").Update(&link)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, link)

	//TODO 重新启动服务

}

func linkGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	err := db.DB("link").One("Id", pid.Id, &link)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, link)
}