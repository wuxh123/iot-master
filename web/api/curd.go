package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"iot-master/db"
	"net/http"
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

func parseBody(request *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

func curdApiList(mod reflect.Type) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		datum := createSliceFromType(mod)

		var body paramSearch
		err := ctx.ShouldBindJSON(&body)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if body.Length < 1 {
			body.Length = 20
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

		for _, keyword := range body.Keywords {
			if keyword.Value != "" {
				op.And(keyword.Key + " like", "%" + keyword.Value + "%")
			}
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
		cnt, err := op.FindAndCount(datum)
		if err != nil {
			replyError(ctx, err)
			return
		}

		//replyOk(ctx, cs)
		replyList(ctx, datum, cnt)
	}
}

func curdApiCreate(mod reflect.Type, before, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := reflect.New(mod).Interface()
		err := ctx.ShouldBindJSON(&data)
		if err != nil {
			replyError(ctx, err)
			return
		}


		if before != nil {
			if err := before(data); err != nil {
				replyError(ctx, err)
				return
			}
		}

		_, err = db.Engine.Insert(data)
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

func curdApiModify(mod reflect.Type, updateFields []string, before, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		err := ctx.ShouldBindUri(&pid)
		if err != nil {
			replyError(ctx, err)
			return
		}

		data := reflect.New(mod).Interface()
		err = ctx.ShouldBindJSON(&data)
		if err != nil {
			replyError(ctx, err)
			return
		}


		if before != nil {
			if err := before(pid.Id); err != nil {
				replyError(ctx, err)
				return
			}
		}

		_, err = db.Engine.ID(pid.Id).Cols(updateFields...).Update(data)
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

func curdApiDelete(mod reflect.Type, before, after hook) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		err := ctx.ShouldBindUri(&pid)
		if err != nil {
			replyError(ctx, err)
			return
		}

		if before != nil {
			if err := before(pid.Id); err != nil {
				replyError(ctx, err)
				return
			}
		}

		data := reflect.New(mod).Interface()
		_, err = db.Engine.ID(pid.Id).Delete(data)
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

func curdApiGet(mod reflect.Type) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pid paramId
		err := ctx.ShouldBindUri(&pid)
		if err != nil {
			replyError(ctx, err)
			return
		}

		data := reflect.New(mod).Interface()
		has, err := db.Engine.ID(pid.Id).Get(data)
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