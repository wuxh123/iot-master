package models

type Template struct {
	UUID        string `json:"uuid"` //唯一码，自动生成
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Links      []TemplateLink `json:"links"`
	Validators []TemplateValidator
	Functions  []TemplateFunction
	Job        []TemplateJob
	Strategies []TemplateStrategy
}

type TemplateLink struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`

	Elements []TemplateElement `json:"elements"`
}

type TemplateElement struct {
	Name  string `json:"name"`
	Alias string `json:"alias"` //别名，用于编程
	Slave uint8  `json:"slave"` //从站号

	Variables []Variable `json:"variables"`
}

type TemplateValidator struct {
	Alert      string `json:"alert"`
	Expression string `json:"expression"` //表达式，检测变量名
}

type TemplateFunction struct {
	Name        string `json:"name"` //项目功能脚本唯一，供外部调用
	Description string `json:"description"`
	Script      string `json:"script"` //javascript
}

type TemplateJob struct {
	Cron     string `json:"cron"`
	Function string `json:"function"`
}

type TemplateStrategy struct {
	Expression string `json:"expression"` //触发条件 表达式，检测变量名
	Function   string `json:"function"`
}
