package api

import (
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type hook func(ctx *gin.Context) error

func curdApiList(typ reflect.Type) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		cs := reflect.MakeSlice(typ, 0, 0)

		var body paramSearch
		err := ctx.ShouldBind(&body)
		if err != nil {
			replyError(ctx, err)
			return
		}

		op := db.Engine.Limit(body.Length, body.Offset)

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
}

func curdApiCreate(typ reflect.Type, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		batch := reflect.New(typ).Interface()
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
}

func curdApiModify(typ reflect.Type, updateFields []string, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		if err := ctx.BindUri(&pid); err != nil {
			replyError(ctx, err)
			return
		}

		batch := reflect.New(typ).Interface()
		if err := ctx.ShouldBindJSON(&batch); err != nil {
			replyError(ctx, err)
			return
		}

		//log.Println("update", batch)
		//TODO 补充列
		_, err := db.Engine.ID(pid.Id).Cols(updateFields...).Update(batch)
		if err != nil {
			replyError(ctx, err)
			return
		}

		replyOk(ctx, batch)
	}
}

func curdApiDelete(typ reflect.Type, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		if err := ctx.BindUri(&pid); err != nil {
			replyError(ctx, err)
			return
		}

		batch := reflect.New(typ).Interface()
		_, err := db.Engine.ID(pid.Id).Delete(batch)
		if err != nil {
			replyError(ctx, err)
			return
		}
		replyOk(ctx, nil)
	}
}

func curdApiGet(typ reflect.Type) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		if err := ctx.BindUri(&pid); err != nil {
			replyError(ctx, err)
			return
		}
		batch := reflect.New(typ).Interface()
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
}
