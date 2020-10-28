package models

import "time"

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"` //模板ID
	UUID        string `json:"uuid"`   //唯一码，自动生成

	Version  string `json:"version"`
	Disabled bool   `json:"disabled"`

	Elements []struct {
		ElementId int    `json:"element_id"`
		LinkId    int    `json:"link_id"`
		Slave     uint8  `json:"slave"` //从站号
		Name      string `json:"name"`
		Alias     string `json:"alias"` //别名，用于编程

		Variables []struct {
			VariableId int     `json:"variable_id"`
			Alias      string  `json:"alias"`   //别名，用于编程
			Correct    float32 `json:"correct"` //校准
		}
	}

	Jobs []struct {
		Name   string `json:"name"`
		Cron   string `json:"cron"`
		Script string `json:"script"` //javascript
	}

	Strategies []struct {
		Name   string `json:"name"`
		Script string `json:"script"` //javascript
	}

	Created time.Time `json:"created" storm:"created"`
	Updated time.Time `json:"updated" storm:"updated"`
}
