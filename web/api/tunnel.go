package api

import (
	"git.zgwit.com/zgwit/MyDTU/core"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func tunnelStart(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		replyError(writer, err)
		return
	}
	c, err := core.GetTunnel(id)
	if err != nil {
		replyError(writer, err)
		return
	}

	err = c.Open()
	if err != nil {
		replyError(writer, err)
		return
	}

	replyOk(writer, nil)
}

func tunnelStop(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		replyError(writer, err)
		return
	}
	c, err := core.GetTunnel(id)
	if err != nil {
		replyError(writer, err)
		return
	}

	err = c.Close()
	if err != nil {
		replyError(writer, err)
		return
	}

	replyOk(writer, nil)
}
