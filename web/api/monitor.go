package api

import (
	"github.com/gin-gonic/gin"
	"git.zgwit.com/iot/dtu-admin/dtu"
)

func monitor(ctx *gin.Context)  {
	var pid paramId2
	if err := ctx.BindUri(&pid); err != nil {
		replyError(ctx, err)
		return
	}

	lnk, err := dtu.GetLink(pid.Id, pid.Id2)
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
		Conn: ws,
		Link: lnk,
	}
	err = lnk.Monitor(m)
	if err != nil {
	}

	//阻塞执行?
	go m.Receive()
}
