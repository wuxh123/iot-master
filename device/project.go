package device

import (
	"iot-master/model"
	"iot-master/types"
)

type project struct {
	model.Project

	//TODO 支持多个链接
	link types.Link

	devices types.Device
}
