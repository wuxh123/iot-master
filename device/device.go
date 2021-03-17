package device

import (
	"errors"
	"iot-master/model"
	"iot-master/types"
	"sync"
)

type device struct {
	model.Device
}


func (d *device)Read(name string) (interface{}, error) {
	//TODO 读变量
	return nil, nil
}
func (d *device)Write(name string, value interface{}) error {
	//TODO 写变量
	return nil
}


var devices sync.Map

func GetDevice(id int64) (types.Device, error) {
	v, ok := projects.Load(id)
	if !ok {
		return nil, errors.New("项目不存在")
	}
	return v.(types.Device), nil
}

