package model

import (
	"git.zgwit.com/zgwit/iot-admin/interfaces"
	"git.zgwit.com/zgwit/iot-admin/types"
	"github.com/robertkrimen/otto"
)

type _link struct {
	linker   interfaces.Linker
	Protocol string
}

type _variable struct {
	Link     *_link
	Type     types.DataType
	Addr     string
	Default  string
	Writable bool //可写，用于输出（如开关）

	//TODO 采样：无、定时、轮询
	Cron string

	Children map[string]_variable
}

type _batchResult struct {
	Offset   int
	Variable string //_variable path
}

type _batch struct {
	Link string
	Type string
	Addr string
	Size int
	Cron string

	Results []_batchResult
}

type _job struct {
	Cron   string
	Script string //javascript
}

type _strategy struct {
	Script string //javascript
}

type Instance struct {
	Links      map[string]_link
	Variables  map[string]_variable
	Batches    map[string]_batch
	Jobs       map[string]_job
	Strategies map[string]_strategy

	vm *otto.Otto
}
