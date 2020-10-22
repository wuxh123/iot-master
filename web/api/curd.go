package api

import (
	"git.zgwit.com/zgwit/iot-admin/db"
	"github.com/kataras/iris/v12"
	"reflect"
)

type hook func(value interface{}) error

func createSliceFromType(mod reflect.Type) interface{} {
	//datas := reflect.MakeSlice(reflect.SliceOf(mod), 0, 10).Interface()

	//解决不可寻址的问题，参考modern-go/reflect2 safe_slice.go
	val := reflect.MakeSlice(reflect.SliceOf(mod), 0, 1)
	ptr := reflect.New(val.Type())
	ptr.Elem().Set(val)
	return ptr.Interface()
}

func curdApiList(mod reflect.Type) iris.Handler {
	return func(ctx iris.Context) {
		datas := createSliceFromType(mod)

		var body paramSearch
		err := ctx.ReadJSON(&body)
		if err != nil {
			replyError(ctx, err)
			return
		}

		op := db.Engine.Limit(body.Length, body.Offset)

		for _, filter := range body.Filters {
			if len(filter.Values) > 0 {
				if len(filter.Values) == 1 {
					op.And(filter.Key+"=?", filter.Values[0])
				} else {
					op.In(filter.Key, filter.Values)
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
		ctx.JSON(iris.Map{
			"ok":    true,
			"data":  datas,
			"total": cnt,
		})
	}
}

func curdApiListById(mod reflect.Type, field string) iris.Handler {
	return func(ctx iris.Context) {
		datas := createSliceFromType(mod)

		id, err := ctx.URLParamInt64("id")
		if err != nil {
			replyError(ctx, err)
			return
		}

		var body paramSearch
		err = ctx.ReadJSON(&body)
		if err != nil {
			replyError(ctx, err)
			return
		}

		op := db.Engine.Where(field+"=?", id).Limit(body.Length, body.Offset)

		for _, filter := range body.Filters {
			if len(filter.Values) > 0 {
				if len(filter.Values) == 1 {
					op.And(filter.Key+"=?", filter.Values[0])
				} else {
					op.In(filter.Key, filter.Values)
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
		ctx.JSON(iris.Map{
			"ok":    true,
			"data":  datas,
			"total": cnt,
		})
	}
}

func curdApiCreate(mod reflect.Type, after hook) iris.Handler {
	return func(ctx iris.Context) {
		data := reflect.New(mod).Interface()
		if err := ctx.ReadJSON(data); err != nil {
			replyError(ctx, err)
			return
		}

		_, err := db.Engine.Insert(data)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(ctx, err)
				return
			}
		}

		replyOk(ctx, data)
	}
}

func curdApiModify(mod reflect.Type, updateFields []string, after hook) iris.Handler {
	return func(ctx iris.Context) {
		id, err := ctx.URLParamInt64("id")
		if err != nil {
			replyError(ctx, err)
			return
		}

		data := reflect.New(mod).Interface()
		if err := ctx.ReadJSON(data); err != nil {
			replyError(ctx, err)
			return
		}

		_, err = db.Engine.ID(id).Cols(updateFields...).Update(data)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(ctx, err)
				return
			}
		}

		replyOk(ctx, data)
	}
}

func curdApiDelete(mod reflect.Type, after hook) iris.Handler {
	return func(ctx iris.Context) {
		id, err := ctx.URLParamInt64("id")
		if err != nil {
			replyError(ctx, err)
			return
		}

		data := reflect.New(mod).Interface()
		_, err = db.Engine.ID(id).Delete(data)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(ctx, err)
				return
			}
		}

		replyOk(ctx, nil)
	}
}

func curdApiGet(mod reflect.Type) iris.Handler {
	return func(ctx iris.Context) {
		id, err := ctx.URLParamInt64("id")
		if err != nil {
			replyError(ctx, err)
			return
		}
		data := reflect.New(mod).Interface()
		has, err := db.Engine.ID(id).Get(data)
		if !has {
			replyFail(ctx, "记录不存在")
			return
		} else if err != nil {
			replyError(ctx, err)
			return
		}
		replyOk(ctx, data)
	}
}
