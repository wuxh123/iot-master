package model

import "time"

type User struct {
	//Id，自增
	Id   int64 `json:"id"`

	//用户名
	Username string  `json:"username" xorm:"varchar(64) notnull unique"`

	//密码 MD5加密
	Password string  `json:"password" xorm:"varchar(64) notnull"`

	//姓名
	Name string `json:"name" xorm:"varchar(64)"`

	//是否禁用
	Disabled bool `json:"disabled,omitempty" xorm:"default 0"`

	//创建时间
	Created time.Time `json:"created" xorm:"created"`
}
