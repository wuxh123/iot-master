package model

import (
	"github.com/robertkrimen/otto"
	"sync"
)

type Variable struct {
	Addr     string
	Type     DataType
	Value    interface{}
	children sync.Map
}

type Simpling struct {
	Cron string
	Type string //read batch
	Path string //path or batch name
}

type Job struct {
	Name   string //唯一
	Cron   string
	Script string
}

type Strategy struct {
	Name      string //唯一
	Script    string //javascript 表达式
	Variables []string
}

type Model struct {
	vm *otto.Otto

	variables sync.Map

	jobs       []Job
	strategies []Strategy
}

func NewModel() *Model {
	return &Model{
		vm:         otto.New(),
		jobs:       make([]Job, 0),
		strategies: make([]Strategy, 0),
	}
}
