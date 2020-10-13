package api

import (
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func batches(ctx *gin.Context) {
	cs := make([]models.ModelBatch, 0)

	var body paramSearch
	err := ctx.ShouldBind(&body)
	if err != nil {
		replyError(ctx, err)
		return
	}

	op := db.Engine.Where("type=?", 0)
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
		op.And("user like ? or text like ? or file like ?", kw, kw, kw)
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

func batchCreate(ctx *gin.Context) {
	var batch models.ModelBatch
	if err := ctx.ShouldBindJSON(&batch); err != nil {
		replyError(ctx, err)
		return
	}

	_, err := db.Engine.Insert(&batch)
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

	_, err := db.Engine.ID(pid.Id).Delete(&models.ModelBatch{})
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

	var batch models.ModelBatch
	if err := ctx.ShouldBindJSON(&batch); err != nil {
		replyError(ctx, err)
		return
	}

	//log.Println("update", batch)
	//TODO 补充列
	_, err := db.Engine.ID(pid.Id).Cols("type", "addr", "size").Update(&models.ModelBatch{})
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
	var batch models.ModelBatch
	has, err := db.Engine.ID(pid.Id).Get(&batch)
	if !has {
		replyFail(ctx, "记录不存在")
		return
	} else if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, batch)
}
