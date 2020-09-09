package api

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/model"
	"github.com/zgwit/storm/v3/q"
	"log"
	"net/http"
)

func links(ctx *gin.Context) {
	var ls []model.Link

	var body paramSearch
	err := ctx.ShouldBind(&body)
	if err != nil {
		replyError(ctx, err)
		return
	}

	cond := make([]q.Matcher, 0)
	//过滤条件
	for _, filter := range body.Filters {
		cond = append(cond, q.In(filter.Key, filter.Value))
	}
	//关键字
	if body.Keyword != "" {
		cond = append(cond, q.Re("name", body.Keyword), q.Re("serial", body.Keyword), q.Re("addr", body.Keyword))
	}

	query := db.DB("channel").Select(cond...)

	//计算总数
	cnt, err := query.Count(&ls)
	if err != nil {
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
		query = query.OrderBy("id").Reverse()
	}

	err = query.Find(&ls)
	if err != nil {
		replyError(ctx, err)
		return
	}

	//replyOk(ctx, cs)
	ctx.JSON(http.StatusOK, gin.H{
		"ok":    true,
		"data":  ls,
		"total": cnt,
	})
}

func linkDelete(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	err := db.DB("link").DeleteStruct(&model.Link{Id: pid.Id})
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, nil)

	//删除服务
	go func() {
		c, err := dtu.GetChannel(link.ChannelId)
		if err != nil {
			log.Println(err)
			return
		}
		l, err := c.GetLink(link.Id)
		if err != nil {
			log.Println(err)
			return
		}
		_ = l.Close()
		//TODO 强制删除连接
	}()

}

func linkModify(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	if err := ctx.ShouldBindJSON(&link); err != nil {
		replyError(ctx, err)
		return
	}

	err := db.DB("link").Update(&link)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, link)

	//TODO 重新启动服务

}

func linkGet(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	err := db.DB("link").One("id", pid.Id, &link)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, link)
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func linkMonitor(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	err := db.DB("link").One("id", pid.Id, &link)
	if err != nil {
		replyError(ctx, err)
		return
	}

	lnk, err := dtu.GetLink(link.ChannelId, link.Id)
	if err != nil {
		replyError(ctx, err)
		return
	}

	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		replyError(ctx, err)
		return
	}

	m := &dtu.Monitor{
		//Key:  "",
		Conn: ws,
		Link: lnk,
	}
	lnk.Monitor(m)
	m.Receive()
	//replyOk(ctx, nil)
}

type linkSendBody struct {
	IsHex bool   `form:"is_hex"`
	Data  string `form:"data"`
}

func linkSend(ctx *gin.Context) {
	var pid paramId
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	var body linkSendBody
	err := ctx.ShouldBind(&body)
	if err != nil {
		replyError(ctx, err)
		return
	}

	var link model.Link
	err = db.DB("link").One("id", pid.Id, &link)
	if err != nil {
		replyError(ctx, err)
		return
	}

	lnk, err := dtu.GetLink(link.ChannelId, link.Id)
	if err != nil {
		replyError(ctx, err)
		return
	}

	b := []byte(body.Data)
	if body.IsHex {
		b, err = hex.DecodeString(body.Data)
		if err != nil {
			replyError(ctx, err)
			return
		}
	}
	_, err = lnk.Send(b)
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}
