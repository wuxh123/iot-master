package api

import (
	"git.zgwit.com/iot/mydtu/core"
	"github.com/gin-gonic/gin"
)

func tunnelStart(c *gin.Context) {
	var pid paramId
	err := c.ShouldBindUri(&pid)
	if err != nil {
		replyError(c, err)
		return
	}
	t, err := core.GetTunnel(pid.Id)
	if err != nil {
		replyError(c, err)
		return
	}

	err = t.Open()
	if err != nil {
		replyError(c, err)
		return
	}

	replyOk(c, nil)
}

func tunnelStop(c *gin.Context) {
	var pid paramId
	err := c.ShouldBindUri(&pid)
	if err != nil {
		replyError(c, err)
		return
	}
	t, err := core.GetTunnel(pid.Id)
	if err != nil {
		replyError(c, err)
		return
	}

	err = t.Close()
	if err != nil {
		replyError(c, err)
		return
	}

	replyOk(c, nil)
}
