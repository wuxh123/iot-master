package core

import (
	"iot-master/base"
	"iot-master/protocol/adapter"
)

type Device struct {
	link    base.Link
	adapter adapter.Adapter
}
