package api

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zgwit/dtu-admin/db"
	"github.com/zgwit/dtu-admin/dtu"
	"github.com/zgwit/dtu-admin/model"
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

	//op := db.Engine.Where("type=?", body.Net)
	op := db.Engine.NewSession()
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
		op.And("name like ? or serial like ? or addr like ?", kw, kw, kw)
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
	cnt, err := op.FindAndCount(&ls)
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
	has, err := db.Engine.ID(pid.Id).Get(&link)
	if err != nil {
		replyError(ctx, err)
		return
	}
	if !has {
		replyFail(ctx, "记录不存在")
		return
	}

	replyOk(ctx, nil)

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

	_, err := db.Engine.ID(pid.Id).Cols("name", "serial", "addr", "channel_id", "plugin_id").Update(&link)
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
	has, err := db.Engine.ID(pid.Id).Get(&link)
	if err != nil {
		replyError(ctx, err)
		return
	}
	if !has {
		replyFail(ctx, "找不到通道")
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
	has, err := db.Engine.ID(pid.Id).Get(&link)
	if err != nil {
		replyError(ctx, err)
		return
	}
	if !has {
		replyFail(ctx, "找不到通道")
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
	has, err := db.Engine.ID(pid.Id).Get(&link)
	if err != nil {
		replyError(ctx, err)
		return
	}
	if !has {
		replyFail(ctx, "找不到通道")
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
