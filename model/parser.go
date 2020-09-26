package model

type _base struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type _link struct {
	_base
	Protocol string `json:"protocol"`
}

type _variable struct {
	_base
	Link     string `json:"link"`
	Type     string `json:"type"`
	Addr     string `json:"addr"`
	Default  string `json:"default"`
	Writable bool   `json:"writable"` //可写，用于输出（如开关）

	Children []_variable `json:"children"`
}

type _batch struct {
	_base
	Link string `json:"link"`
	Type string `json:"type"`
	Addr string `json:"addr"`
	Size int    `json:"size"`

	Results []struct {
		Offset   int    `json:"offset"`
		Variable string `json:"variable"` //Variable
	} `json:"results"`
}

type _sampling struct {
	_base
	Cron string `json:"cron"`
	Type string `json:"type"` //read batch
	Path string `json:"path"` //path or batch name
}

type _job struct {
	_base
	Cron   string `json:"cron"`
	Type   string `json:"type"`   //read batch strategy
	Target string `json:"target"` //path, name, name
}

type _strategy struct {
	_base
	Script string `json:"script"` //javascript 表达式
}

type _model struct {
	_base

	Links      []_link     `json:"links"`
	Variables  []_variable `json:"variables"`
	Batches    []_batch    `json:"batches"`
	Samplings  []_sampling `json:"samplings"`
	Jobs       []_job      `json:"jobs"`
	Strategies []_strategy `json:"strategies"`
}
