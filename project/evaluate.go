package project

import (
	"github.com/Knetic/govaluate"
	"sync"
)

var expressions sync.Map

type Vars []string

func Evaluate(expression string, parameters map[string]interface{}) (interface{}, error) {
	//可以判断长度，过长的表达式不绑在
	var expr *govaluate.EvaluableExpression
	var err error
	//从缓存中加载表达式
	e, ok := expressions.Load(expression)
	if !ok {
		expr, err = govaluate.NewEvaluableExpression(expression)
		if err != nil {
			return nil, err
		}
	} else {
		expr = e.(*govaluate.EvaluableExpression)
	}
	return expr.Evaluate(parameters)
}


func ParseExpressionVariables(expression string) ([]string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return nil, err
	}
	expressions.Store(expression, expr)
	return expr.Vars(), nil
}
