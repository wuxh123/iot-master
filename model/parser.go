package model

type _base struct {
	Name     string
	Description string
}

type _link struct {
	_base
	Protocol string
}

type _variable struct {
	_base
	Link     string
	Type     string
	Addr     string
	Default  string
	Writable bool //可写，用于输出（如开关）

	Children []_variable
}

type _batch struct {
	_base
	Link string
	Type string
	Addr string
	Size int

	Results []struct {
		Offset int
		Path   string //Variable
	}
}

type _sampling struct {
	_base
	Cron string
	Type string //read batch
	Path string //path or batch name
}

type _job struct {
	_base
	Cron   string
	Script string //javascript 表达式
}

type _strategy struct {
	_base
	Script string //javascript 表达式
}

type _model struct {
	_base
	
	links      []_link
	variables  []_variable
	batches    []_batch
	samplings  []_sampling
	jobs       []_job
	strategies []_strategy
}
