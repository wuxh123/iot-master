package conf

import (
	"github.com/zgwit/dtu-admin/flag"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type _database struct {
	Desc  string `json:"desc" yaml:"desc"`
	Type  string `json:"type" yaml:"type"`
	Url   string `json:"url" yaml:"url"`
	Debug bool   `json:"showSQL" yaml:"showSQL"`
}

type _web struct {
	Desc  string `yaml:"desc"`
	Addr  string `yaml:"addr"`
	Cors  bool   `yaml:"cors"`
	Debug bool   `yaml:"debug"`
}

type _baseAuth struct {
	Desc     string `yaml:"desc"`
	Disabled bool   `yaml:"disabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type _sysAdmin struct {
	Desc   string `yaml:"desc"`
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"`
}

type _dbus struct {
	Desc string `yaml:"desc"`
	Addr string `yaml:"addr"`
}

type _config struct {
	Database _database `yaml:"database"`
	Web      _web      `yaml:"web"`
	BaseAuth _baseAuth `yaml:"basic_auth"`
	SysAdmin _sysAdmin `yaml:"sys_admin"`
	DBus     _dbus     `yaml:"dbus"`
}

var Config = _config{
	Database: _database{
		Desc: "数据库配置",
		Type: "mysql",
		Url:  "root:root@tcp(127.0.0.1:3306)/dtu-admin?charset=utf8",
	},
	Web: _web{
		Desc: "Web服务配置",
		Addr: ":8080",
	},
	BaseAuth: _baseAuth{
		Desc:     "HTTP简单认证，仅用于超级管理员",
		Username: "admin",
		Password: "123456",
	},
	SysAdmin: _sysAdmin{
		Desc: "Sys Admin地址",
		Addr: "http://127.0.0.1:8080",
	},
	DBus: _dbus{
		Desc: "数据总线",
		Addr: ":1843",
	},
}

func Load() {
	log.Println("加载配置")

	//从参数中读取配置文件名
	filename := flag.ConfigPath

	// 如果没有文件，则使用默认信息创建
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		Save()
	} else {
		y, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer y.Close()

		d := yaml.NewDecoder(y)
		_ = d.Decode(&Config)
	}
}

func Save() {
	log.Println("保存配置")
	//从参数中读取配置文件名
	filename := flag.ConfigPath

	y, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755) //os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer y.Close()

	e := yaml.NewEncoder(y)
	defer e.Close()

	_ = e.Encode(&Config)
}
