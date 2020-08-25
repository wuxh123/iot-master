package conf

import (
	"github.com/zgwit/dtu-admin/flag"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type _database struct {
	Desc    string `yaml:"desc"`
	Type    string `yaml:"type"`
	Url     string `yaml:"url"`
	ShowSQL bool   `yaml:"showSQL"`
}

type _http struct {
	Desc  string `yaml:"desc"`
	Addr  string `yaml:"addr"`
	Debug bool   `yaml:"debug"`
	Cors  bool   `yaml:"cors"`
}

type _config struct {
	Database _database `yaml:"database"`
	Http     _http     `yaml:"http"`
}

var Config = _config{
	_database{
		"数据库配置",
		"mysql",
		"root:root@tcp(127.0.0.1:3306)/dtu?charset=utf8",
		false,
	},
	_http{
		"HTTP服务配置，API接口",
		":8080",
		false,
		false,
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
