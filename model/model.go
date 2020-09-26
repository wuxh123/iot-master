package model

import (
	"github.com/robertkrimen/otto"
)

type Variable struct {
	Name     string //唯一
	Addr     string
	Type     DataType
	Value    interface{}
	Default  interface{}
	Writable bool //可写，用于输出（如开关）
	Children map[string]Variable
}

type Simpling struct {
	Name string //唯一
	Cron string
	Type string //read batch
	Path string //path or batch name
}

type Job struct {
	Name   string //唯一
	Cron   string
	Script string //javascript 表达式
}

type Strategy struct {
	Name   string //唯一
	Script string //javascript 表达式
}

type Model struct {
	vm *otto.Otto

	variables  []Variable
	strategies []Strategy
	jobs       []Job
}

func NewModel() *Model {
	return &Model{
		vm: otto.New(),
	}
}
