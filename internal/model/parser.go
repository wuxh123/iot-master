package model

type _base struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModelLink struct {
	_base
	Protocol string `json:"protocol"`
}

type ModelVariable struct {
	_base
	Link     string `json:"link"`
	Type     string `json:"type"`
	Addr     string `json:"addr"`
	Default  string `json:"default"`
	Writable bool   `json:"writable"` //可写，用于输出（如开关）

	//TODO 采样：无、定时、轮询
	Cron string `json:"cron"`

	Children []ModelVariable `json:"children"`
}

type ModelBatchResult struct {
	Offset   int    `json:"offset"`
	Variable string `json:"variable"` //ModelVariable path
}

type ModelBatch struct {
	_base
	Link string `json:"link"`
	Type string `json:"type"`
	Addr string `json:"addr"`
	Size int    `json:"size"`
	Cron string `json:"cron"`

	Results []ModelBatchResult `json:"results"`
}

type ModelJob struct {
	_base
	Cron   string `json:"cron"`
	Script string `json:"script"` //javascript
}

type ModelStrategy struct {
	_base
	Script string `json:"script"` //javascript
}

type Model struct {
	_base

	Version string `json:"version"`
	H5      string `json:"h5"`

	Links      []ModelLink     `json:"links"`
	Variables  []ModelVariable `json:"variables"`
	Batches    []ModelBatch    `json:"batches"`
	Jobs       []ModelJob      `json:"jobs"`
	Strategies []ModelStrategy `json:"strategies"`
}

func Import(json string) error {
	//TODO parser model, import
	return nil
}

func Export(id int) string {
	//TODO ge
	return ""
}
