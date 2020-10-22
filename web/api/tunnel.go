package api

import (
	"git.zgwit.com/zgwit/iot-admin/core"
	"github.com/kataras/iris/v12"
)

func tunnelStart(ctx iris.Context) {
	id, err := ctx.URLParamInt64("id")
	if err != nil {
		replyError(ctx, err)
		return
	}
	c, err := core.GetTunnel(id)
	if err != nil {
		replyError(ctx, err)
		return
	}

	err = c.Open()
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}

func tunnelStop(ctx iris.Context) {
	id, err := ctx.URLParamInt64("id")
	if err != nil {
		replyError(ctx, err)
		return
	}
	c, err := core.GetTunnel(id)
	if err != nil {
		replyError(ctx, err)
		return
	}

	err = c.Close()
	if err != nil {
		replyError(ctx, err)
		return
	}

	replyOk(ctx, nil)
}
