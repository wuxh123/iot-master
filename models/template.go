package models

type Template struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Jobs       []TemplateJob      `json:"jobs"`
	Strategies []TemplateStrategy `json:"strategies"`
}

type TemplateJob struct {
	Name   string `json:"name"`
	Cron   string `json:"cron"`
	Script string `json:"script"` //javascript
}

type TemplateStrategy struct {
	Name   string `json:"name"`
	Script string `json:"script"` //javascript
}
