package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
	"io/ioutil"
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

//type Handler func(c *gin.Context)

func curdApiList(store storm.Node, mod reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		datas := createSliceFromType(mod)
		data := reflect.New(mod).Interface()

		var body paramSearch
		err := c.ShouldBindJSON(&body)
		if err != nil {
			replyError(c, err)
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
			replyError(c, err)
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
			replyError(c, err)
			return
		}

		//replyOk(ctx, cs)
		replyList(c, datas, cnt)
	}
}

func curdApiListById(store storm.Node,mod reflect.Type, field string) gin.HandlerFunc {
	return func(c *gin.Context) {
		datas := createSliceFromType(mod)
		data := reflect.New(mod).Interface()

		var pid paramId
		err := c.ShouldBindUri(&pid)
		if err != nil {
			replyError(c, err)
			return
		}

		var body paramSearch
		err = c.ShouldBindJSON(&body)
		if err != nil {
			replyError(c, err)
			return
		}

		cond := make([]q.Matcher, 0)
		cond = append(cond, q.Eq(field, pid.Id))

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
			replyError(c, err)
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
			replyError(c, err)
			return
		}

		//replyOk(ctx, cs)
		replyList(c, datas, cnt)
	}
}

func curdApiCreate(store storm.Node,mod reflect.Type, before hook, after hook) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := reflect.New(mod).Interface()
		if err := c.ShouldBindJSON(data); err != nil {
			replyError(c, err)
			return
		}

		if before != nil {
			if err := before(data); err != nil {
				replyError(c, err)
				return
			}
		}

		err := store.Save(data)
		if err != nil {
			replyError(c, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(c, err)
				return
			}
		}

		replyOk(c, data)
	}
}

func curdApiModify(store storm.Node,mod reflect.Type, before hook, after hook) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pid paramId
		err := c.ShouldBindUri(&pid)
		if err != nil {
			replyError(c, err)
			return
		}

		val := reflect.New(mod)
		data := val.Interface()
		if err := c.ShouldBindJSON(data); err != nil {
			replyError(c, err)
			return
		}

		val.Elem().FieldByName("ID").Set(reflect.ValueOf(pid.Id))

		if before != nil {
			if err := before(data); err != nil {
				replyError(c, err)
				return
			}
		}

		err = store.Update(data)
		if err != nil {
			replyError(c, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(c, err)
				return
			}
		}

		replyOk(c, data)
	}
}

func curdApiDelete(store storm.Node,mod reflect.Type, before hook, after hook) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pid paramId
		err := c.ShouldBindUri(&pid)
		if err != nil {
			replyError(c, err)
			return
		}

		val := reflect.New(mod)
		data := val.Interface()
		val.Elem().FieldByName("ID").Set(reflect.ValueOf(pid.Id))


		if before != nil {
			if err := before(pid.Id); err != nil {
				replyError(c, err)
				return
			}
		}

		err = store.DeleteStruct(data)
		if err != nil {
			replyError(c, err)
			return
		}

		if after != nil {
			err = after(data)
			if err != nil {
				replyError(c, err)
				return
			}
		}

		replyOk(c, nil)
	}
}

func curdApiGet(store storm.Node,mod reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pid paramId
		err := c.ShouldBindUri(&pid)
		if err != nil {
			replyError(c, err)
			return
		}
		data := reflect.New(mod).Interface()
		err = store.One("ID", pid.Id, data)
		if err != nil {
			replyError(c, err)
			return
		}
		replyOk(c, data)
	}
}
