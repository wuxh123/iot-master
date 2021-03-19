package modbus

import "iot-master/protocol"

var manifest = protocol.Manifest{
	Name:    "Modbus RTU",
	Version: "1.0",
	Codes: []protocol.Code{
		{"线圈", 1},
		{"离散量", 2},
		{"保持寄存器", 3},
		{"输入寄存器", 4},
	},
}

func Manifest() protocol.Manifest {
	return manifest
}


func init() {

	protocol.RegisterProtocol(
		"Modbus RTU",
		&manifest,
		NewModbusRtu)

}
