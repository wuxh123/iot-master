package api


import (
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/types"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
	"net/http"
)

func batches(ctx *gin.Context) {
	cs := make([]types.ModelBatch, 0)

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

	query := db.DB("model").From("batch").Select(cond...)

	//计算总数
	cnt, err := query.Count(&types.ModelBatch{})
	if err != nil && err != storm.ErrNotFound {
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
	if err != nil && err != storm.ErrNotFound {
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

func batchCreate(ctx *gin.Context) {
	var batch types.ModelBatch
	if err := ctx.ShouldBindJSON(&batch); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("model").From("batch").Save(&batch)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, batch)
}

func batchDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("model").From("batch").DeleteStruct(&types.Link{Id: pid.Id})
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, nil)
}

func batchModify(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var batch types.ModelBatch
	if err := ctx.ShouldBindJSON(&batch); err != nil {
		replyError(ctx, err)
		return
	}

	//log.Println("update", batch)
	err := db.DB("model").From("batch").Update(&batch)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, batch)
}


func batchGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}
	var batch types.ModelBatch
	err := db.DB("model").From("batch").One("Id", pid.Id, &batch)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, batch)
}

