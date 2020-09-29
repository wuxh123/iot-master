package model

import (
	"github.com/robertkrimen/otto"
)


type Instance struct {
	//Links      map[string]_link
	Variables  map[string]_variable
	//Batches    map[string]_batch
	//Jobs       map[string]_job
	//Strategies map[string]_strategy

	vm *otto.Otto
}


