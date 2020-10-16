package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
//"github.com/modern-go/reflect2"
)

type hook func(value interface{}) error

func curdApiList(mod reflect.Type) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//datas := reflect.MakeSlice(reflect.SliceOf(mod), 0, 10).Interface()

		//解决不可寻址的问题，参考modern-go/reflect2 safe_slice.go
		val := reflect.MakeSlice(reflect.SliceOf(mod), 0, 1)
		ptr := reflect.New(val.Type())
		ptr.Elem().Set(val)
		datas := ptr.Interface()

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
		cnt, err := op.FindAndCount(datas)
		if err != nil {
			replyError(ctx, err)
			return
		}

		//replyOk(ctx, cs)
		ctx.JSON(http.StatusOK, gin.H{
			"ok":    true,
			"data":  datas,
			"total": cnt,
		})
	}
}

func curdApiCreate(mod reflect.Type, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		batch := reflect.New(mod).Interface()
		if err := ctx.ShouldBindJSON(&batch); err != nil {
			replyError(ctx, err)
			return
		}

		_, err := db.Engine.Insert(&batch)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			err = after(batch)
			if err != nil {
				replyError(ctx, err)
				return
			}
		}

		replyOk(ctx, batch)
	}
}

func curdApiModify(mod reflect.Type, updateFields []string, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		if err := ctx.BindUri(&pid); err != nil {
			replyError(ctx, err)
			return
		}

		batch := reflect.New(mod).Interface()
		if err := ctx.ShouldBindJSON(&batch); err != nil {
			replyError(ctx, err)
			return
		}

		_, err := db.Engine.ID(pid.Id).Cols(updateFields...).Update(batch)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			err = after(batch)
			if err != nil {
				replyError(ctx, err)
				return
			}
		}

		replyOk(ctx, batch)
	}
}

func curdApiDelete(mod reflect.Type, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		if err := ctx.BindUri(&pid); err != nil {
			replyError(ctx, err)
			return
		}

		batch := reflect.New(mod).Interface()
		_, err := db.Engine.ID(pid.Id).Delete(batch)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			err = after(batch)
			if err != nil {
				replyError(ctx, err)
				return
			}
		}

		replyOk(ctx, nil)
	}
}

func curdApiGet(mod reflect.Type) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		if err := ctx.BindUri(&pid); err != nil {
			replyError(ctx, err)
			return
		}
		batch := reflect.New(mod).Interface()
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
