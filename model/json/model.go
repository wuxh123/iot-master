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
	Size int

	Results []struct {
		Offset int
		Path   string //Variable
	}
}

type Simpling struct {
	Cron string
	Type string //read batch
	Path string
}

type Job struct {
	Cron   string
	Script string
}

type Strategy struct {
	Variables []string
	Script    string //javascript 表达式
}
