package device

import (
	"errors"
	"iot-master/model"
	"iot-master/types"
	"sync"
)

type project struct {
	model.Project

	//支持多个链接
	links []types.Link
}

func (p *project)Run(name string) error  {
	//TODO 调用脚本

	return nil
}


var projects sync.Map

func GetProject(id int64) (types.Project, error) {
	v, ok := projects.Load(id)
	if !ok {
		return nil, errors.New("项目不存在")
	}
	return v.(types.Project), nil
}
