package project

import (
	"iot-master/model"
)

type device struct {
	model.Device

	variables map[string]interface{}

	validatorsVars []Vars
}

func (d *device) Read(name string) (interface{}, error) {
	//TODO 读变量
	return nil, nil
}
func (d *device) Write(name string, value interface{}) error {
	//TODO 写变量
	return nil
}
