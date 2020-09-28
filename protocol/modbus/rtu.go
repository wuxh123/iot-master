package modbus

import (
	"encoding/binary"
	"git.zgwit.com/zgwit/iot-admin/interfaces"
	"git.zgwit.com/zgwit/iot-admin/protocol"
	"git.zgwit.com/zgwit/iot-admin/protocol/helper"
)

type RTU struct {
	linker   interfaces.Linker
	listener protocol.AdapterListener

	addr *address
}

func NewModbusRtu(linker interfaces.Linker) protocol.Adapter {
	m := &RTU{linker: linker}
	linker.Listen(m)
	//TODO un listen
	return m
}

func (m *RTU) Listen(listener protocol.AdapterListener) {
	m.listener = listener
}

func (m *RTU) Name() string {
	return "Modbus RTU"
}

func (m *RTU) Version() string {
	return "v0.0.1"
}

func (m *RTU) Read(addr string, size int) (err error) {
	m.addr, err = parseAddress(addr)
	if err != nil {
		return err
	}

	// 打包协议
	b := make([]byte, 8)
	b[0] = m.addr.unit
	b[1] = m.addr.code
	helper.WriteUint16(b[2:], m.addr.addr)
	helper.WriteUint16(b[4:], uint16(size))
	helper.WriteUint16(b[6:], CRC16(b[:6]))

	return m.linker.Write(b)
}

func (m *RTU) Write(addr string, buf []byte) (err error) {
	m.addr, err = parseAddress(addr)
	if err != nil {
		return err
	}

	// 打包协议
	l := 6 + len(buf)
	b := make([]byte, l)
	b[0] = m.addr.unit
	b[1] = m.addr.code
	helper.WriteUint16(b[2:], m.addr.addr)
	copy(b[4:], buf)
	helper.WriteUint16(b[l-2:], CRC16(b[:l-2]))

	return m.linker.Write(b)
}

func (m *RTU) OnLinkerData(buf []byte) {
	//TODO 解析数据

	l := len(buf)
	crc := helper.ParseUint16(buf[l-2:])

	if crc != CRC16(buf[:l-2]) {
		//检验错误
		return
	}

	length := 4
	switch buf[1] {
	case FuncCodeReadDiscreteInputs,
		FuncCodeReadCoils:
		count := int(helper.ParseUint16(buf[4:]))
		length += 1 + count/8
		if count%8 != 0 {
			length++
		}
	case FuncCodeReadInputRegisters,
		FuncCodeReadHoldingRegisters,
		FuncCodeReadWriteMultipleRegisters:
		count := int(helper.ParseUint16(buf[4:]))
		length += 1 + count*2
	case FuncCodeWriteSingleCoil,
		FuncCodeWriteMultipleCoils,
		FuncCodeWriteSingleRegister,
		FuncCodeWriteMultipleRegisters:
		length += 4
	case FuncCodeMaskWriteRegister:
		length += 6
	case FuncCodeReadFIFOQueue:
		// undetermined
	default:

	}

	if l < length {
		//长度不够
		return
	}
	b := buf[2 : l-2]

	//m.listener.OnAdapterRead()
	//m.listener.OnAdapterWrite()

}

func (m *RTU) OnLinkerError(err error) {
	m.listener.OnAdapterError(err)
}

func (m *RTU) OnLinkerClose() {}
