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

	//TODO 采样：无、定时、轮询
	Cron string `json:"cron"`

	Children []_variable `json:"children"`
}

type _batch struct {
	_base
	Link string `json:"link"`
	Type string `json:"type"`
	Addr string `json:"addr"`
	Size int    `json:"size"`
	Cron string `json:"cron"`

	Results []struct {
		Offset   int    `json:"offset"`
		Variable string `json:"variable"` //Variable path
	} `json:"results"`
}

type _job struct {
	_base
	Cron   string `json:"cron"`
	Script string `json:"script"` //javascript
}

type _strategy struct {
	_base
	Script string `json:"script"` //javascript
}

type _model struct {
	_base

	Links      []_link     `json:"links"`
	Variables  []_variable `json:"variables"`
	Batches    []_batch    `json:"batches"`
	Jobs       []_job      `json:"jobs"`
	Strategies []_strategy `json:"strategies"`
}
