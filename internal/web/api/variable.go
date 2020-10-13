package api


import (
	"git.zgwit.com/zgwit/iot-admin/internal/db"
	"git.zgwit.com/zgwit/iot-admin/models"
	"github.com/gin-gonic/gin"
	"github.com/zgwit/storm/v3"
	"github.com/zgwit/storm/v3/q"
	"net/http"
)

func variables(ctx *gin.Context) {
	cs := make([]models.ModelVariable, 0)

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

	query := db.DB("model").From("variable").Select(cond...)

	//计算总数
	cnt, err := query.Count(&models.ModelVariable{})
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

func variableCreate(ctx *gin.Context) {
	var variable models.ModelVariable
	if err := ctx.ShouldBindJSON(&variable); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("model").From("variable").Save(&variable)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, variable)
}

func variableDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("model").From("variable").DeleteStruct(&models.Link{Id: pid.Id})
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, nil)
}

func variableModify(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var variable models.ModelVariable
	if err := ctx.ShouldBindJSON(&variable); err != nil {
		replyError(ctx, err)
		return
	}

	//log.Println("update", variable)
	err := db.DB("model").From("variable").Update(&variable)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, variable)
}


func variableGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}
	var variable models.ModelVariable
	err := db.DB("model").From("variable").One("Id", pid.Id, &variable)
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, variable)
}

