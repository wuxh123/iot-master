package json

/**
目录结构：/data/models/{model}/...
h5/...
js/...
links.json
variables.json
batches.json
jobs.json


*/

type Link struct {
	Name     string
	Protocol string
}

type Variable struct {
	Link     string
	Type     string
	Addr     string
	Path     string
	Writable bool //可写，用于输出（如开关）
}

type Batch struct {
	Name string //唯一
	Link string
	Type string
	Addr string

	Results []struct {
		Offset int
		Path   string //Variable
	}
}

type Job struct {
	Cron string
	Type string //read write javascript
	Path string //变量 或 脚本
}

type Strategy struct {
	Cond       string //javascript 表达式
	Operations []struct {
		Path  string
		Value string
	}
}
