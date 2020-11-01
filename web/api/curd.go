package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
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

func parseBody(request *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

type Handler func(writer http.ResponseWriter, request *http.Request)

func curdApiList(store storm.Node, mod reflect.Type) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		datas := createSliceFromType(mod)
		data := reflect.New(mod).Interface()

		var body paramSearch
		err := parseBody(request, &body)
		if err != nil {
			replyError(writer, err)
			return
		}

		cond := make([]q.Matcher, 0)

		//过滤
		for _, filter := range body.Filters {
			if len(filter.Values) > 0 {
				if len(filter.Values) == 1 {
					cond = append(cond, q.Eq(filter.Key, filter.Values[0]))
				} else {
					cond = append(cond, q.In(filter.Key, filter.Values))
				}
			}
		}

		//关键字搜索
		kws := make([]q.Matcher, 0)
		for _, keyword := range body.Keywords {
			if keyword.Value != "" {
				kws = append(kws, q.Re(keyword.Key, keyword.Value))
			}
		}
		if len(kws) > 0 {
			cond = append(cond, q.Or(kws...))
		}

		query := store.Select(cond...)

		//计算总数
		cnt, err := query.Count(data)
		if err != nil && err != storm.ErrNotFound {
			replyError(writer, err)
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
			query = query.OrderBy("ID").Reverse()
		}

		err = query.Find(datas)
		if err != nil && err != storm.ErrNotFound {
			replyError(writer, err)
			return
		}

		//replyOk(ctx, cs)
		replyList(writer, datas, cnt)
	}
}

func curdApiListById(store storm.Node,mod reflect.Type, field string) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		datas := createSliceFromType(mod)
		data := reflect.New(mod).Interface()

		id, err := strconv.Atoi(mux.Vars(request)["id"])
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

		cond := make([]q.Matcher, 0)
		cond = append(cond, q.Eq(field, id))

		//过滤
		for _, filter := range body.Filters {
			if len(filter.Values) > 0 {
				if len(filter.Values) == 1 {
					cond = append(cond, q.Eq(filter.Key, filter.Values[0]))
				} else {
					cond = append(cond, q.In(filter.Key, filter.Values))
				}
			}
		}

		//关键字搜索
		kws := make([]q.Matcher, 0)
		for _, keyword := range body.Keywords {
			kws = append(kws, q.Re(keyword.Key, keyword.Value))
		}
		if len(kws) > 0 {
			cond = append(cond, q.Or(kws...))
		}

		query := store.Select(cond...)

		//计算总数
		cnt, err := query.Count(data)
		if err != nil && err != storm.ErrNotFound {
			replyError(writer, err)
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
			query = query.OrderBy("ID").Reverse()
		}

		err = query.Find(datas)
		if err != nil && err != storm.ErrNotFound {
			replyError(writer, err)
			return
		}

		//replyOk(ctx, cs)
		replyList(writer, datas, cnt)
	}
}

func curdApiCreate(store storm.Node,mod reflect.Type, before hook, after hook) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		data := reflect.New(mod).Interface()
		if err := parseBody(request, data); err != nil {
			replyError(writer, err)
			return
		}

		if before != nil {
			if err := before(data); err != nil {
				replyError(writer, err)
				return
			}
		}

		err := store.Save(data)
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

func curdApiModify(store storm.Node,mod reflect.Type, before hook, after hook) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(mux.Vars(request)["id"])
		if err != nil {
			replyError(writer, err)
			return
		}

		val := reflect.New(mod)
		data := val.Interface()
		if err := parseBody(request, data); err != nil {
			replyError(writer, err)
			return
		}

		val.Elem().FieldByName("ID").Set(reflect.ValueOf(id))

		if before != nil {
			if err := before(data); err != nil {
				replyError(writer, err)
				return
			}
		}

		err = store.Update(data)
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

func curdApiDelete(store storm.Node,mod reflect.Type, before hook, after hook) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(mux.Vars(request)["id"])
		if err != nil {
			replyError(writer, err)
			return
		}

		val := reflect.New(mod)
		data := val.Interface()
		val.Elem().FieldByName("ID").Set(reflect.ValueOf(id))


		if before != nil {
			if err := before(id); err != nil {
				replyError(writer, err)
				return
			}
		}

		err = store.DeleteStruct(data)
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

func curdApiGet(store storm.Node,mod reflect.Type) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(mux.Vars(request)["id"])
		if err != nil {
			replyError(writer, err)
			return
		}
		data := reflect.New(mod).Interface()
		err = store.One("ID", id, data)
		if err != nil {
			replyError(writer, err)
			return
		}
		replyOk(writer, data)
	}
}
