package api

import (
	"encoding/json"
	"git.zgwit.com/zgwit/iot-admin/db"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
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

func parseBody(request *http.Request, data interface{}) error  {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

type Handler func(writer http.ResponseWriter, request *http.Request)

func curdApiList(mod reflect.Type) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		datas := createSliceFromType(mod)

		var body paramSearch
		err := parseBody(request, &body)
		if err != nil {
			replyError(writer, err)
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
			replyError(writer, err)
			return
		}

		//replyOk(ctx, cs)
		replyList(writer, datas, cnt)
	}
}

func curdApiListById(mod reflect.Type, field string) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		datas := createSliceFromType(mod)

		id, err := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
		if err != nil {
			replyError(writer, err)
			return
		}

		var body paramSearch
		err = parseBody(request, &body)
		if err != nil {
			replyError(writer, err)
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
			replyError(writer, err)
			return
		}

		//replyOk(ctx, cs)
		replyList(writer, datas, cnt)
	}
}

func curdApiCreate(mod reflect.Type, after hook) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		data := reflect.New(mod).Interface()
		if err := parseBody(request, data); err != nil {
			replyError(writer, err)
			return
		}

		_, err := db.Engine.Insert(data)
		if err != nil {
			replyError(writer, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(writer, err)
				return
			}
		}

		replyOk(writer, data)
	}
}

func curdApiModify(mod reflect.Type, updateFields []string, after hook) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
		if err != nil {
			replyError(writer, err)
			return
		}

		data := reflect.New(mod).Interface()
		if err := parseBody(request, data); err != nil {
			replyError(writer, err)
			return
		}

		_, err = db.Engine.ID(id).Cols(updateFields...).Update(data)
		if err != nil {
			replyError(writer, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(writer, err)
				return
			}
		}

		replyOk(writer, data)
	}
}

func curdApiDelete(mod reflect.Type, after hook) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
		if err != nil {
			replyError(writer, err)
			return
		}

		data := reflect.New(mod).Interface()
		_, err = db.Engine.ID(id).Delete(data)
		if err != nil {
			replyError(writer, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(writer, err)
				return
			}
		}

		replyOk(writer, nil)
	}
}

func curdApiGet(mod reflect.Type) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
		if err != nil {
			replyError(writer, err)
			return
		}
		data := reflect.New(mod).Interface()
		has, err := db.Engine.ID(id).Get(data)
		if !has {
			replyFail(writer, "记录不存在")
			return
		} else if err != nil {
			replyError(writer, err)
			return
		}
		replyOk(writer, data)
	}
}
