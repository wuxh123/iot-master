package modbus

import "iot-master/protocol"

var manifest = protocol.Manifest{
	Name:    "",
	Version: "",
	Codes:   nil,
}

func Manifest() protocol.Manifest  {
	return manifest
}
