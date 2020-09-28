package protocol

import (
	"git.zgwit.com/zgwit/iot-admin/interfaces"
	"git.zgwit.com/zgwit/iot-admin/types"
)

type ModbusRtu struct {
	linker   interfaces.Linker
	listener AdapterListener
}

func NewModbusRtu(linker interfaces.Linker) Adapter {
	m := &ModbusRtu{linker: linker}
	linker.Listen(m)
	//TODO un listen
	return m
}

func (m *ModbusRtu) Listen(listener AdapterListener) {
	m.listener = listener
}

func (m *ModbusRtu) Name() string {
	return "Modbus RTU"
}

func (m *ModbusRtu) Version() string {
	return "v0.0.1"
}

func (m *ModbusRtu) Read(addr string, typ types.DataType, size int) (err error) {
	//TODO 打包协议
	return m.linker.Write([]byte(""))
}

func (m *ModbusRtu) Write(addr string, typ types.DataType, buf []byte) (err error) {
	//TODO 打包协议
	return m.linker.Write([]byte(""))
}

func (m *ModbusRtu) OnLinkerData(buf []byte) {
	//TODO 解析数据
	//m.listener.OnAdapterRead()
	//m.listener.OnAdapterWrite()
}

func (m *ModbusRtu) OnLinkerError(err error) {
	m.listener.OnAdapterError(err)
}

func (m *ModbusRtu) OnLinkerClose() {}
