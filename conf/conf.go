package conf

import (
	"mydtu/args"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type _data struct {
	Desc  string `yaml:"desc"`
	Path  string `yaml:"path"`
	Debug bool   `yaml:"debug"`
}


type _web struct {
	Desc string `yaml:"desc"`
	Addr string `yaml:"addr"`
	//Cors      bool   `yaml:"cors"`
	Debug bool `yaml:"debug"`
}

type _dbus struct {
	Desc string `yaml:"desc"`
	Addr string `yaml:"addr"`
}

type _config struct {
	Data     _data     `yaml:"data"`
	Web      _web      `yaml:"web"`
	DBus     _dbus     `yaml:"dbus"`
}

var Config = _config{
	Data: _data{
		Desc:    "数据库配置",
		Path:    "data",
	},
	Web: _web{
		Desc: "Web服务配置",
		Addr: ":8080",
	},
	DBus: _dbus{
		Desc: "数据总线",
		Addr: ":1843",
	},
}

func Load() error {
	//log.Println("加载配置")
	//从参数中读取配置文件名
	filename := args.ConfigPath

	// 如果没有文件，则使用默认信息创建
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return Save()
	} else {
		y, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer y.Close()

		d := yaml.NewDecoder(y)
		return d.Decode(&Config)
	}
	return nil
}

func Save() error {
	//log.Println("保存配置")
	//从参数中读取配置文件名
	filename := args.ConfigPath

	y, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755) //os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer y.Close()

	e := yaml.NewEncoder(y)
	defer e.Close()

	return e.Encode(&Config)
}
