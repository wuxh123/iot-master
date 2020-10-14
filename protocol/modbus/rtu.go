package modbus

import (
	"git.zgwit.com/zgwit/iot-admin/base"
	"git.zgwit.com/zgwit/iot-admin/models"
	"git.zgwit.com/zgwit/iot-admin/protocol"
	"git.zgwit.com/zgwit/iot-admin/protocol/helper"
)

type RTU struct {
	link base.Link

	listener protocol.AdapterListener
	//addr *address
}

func NewModbusRtu(linker base.Link) protocol.Adapter {
	m := &RTU{link: linker}
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

func (m *RTU) Read(addr *models.Address, size int) (err error) {
	b := make([]byte, 8)
	b[0] = uint8(addr.Slave)
	b[1] = uint8(addr.ReadCode)
	helper.WriteUint16(b[2:], uint16(addr.Offset))
	helper.WriteUint16(b[4:], uint16(size))
	helper.WriteUint16(b[6:], CRC16(b[:6]))

	return m.link.Write(b)
}

func (m *RTU) Write(addr *models.Address, buf []byte) (err error) {
	//TODO 如果是线圈，需要Shrink
	l := 6 + len(buf)
	b := make([]byte, l)
	b[0] = uint8(addr.Slave)
	b[1] = uint8(addr.WriteCode)
	helper.WriteUint16(b[2:], uint16(addr.Offset))
	copy(b[4:], buf)
	helper.WriteUint16(b[l-2:], CRC16(b[:l-2]))

	return m.link.Write(b)
}

func (m *RTU) OnLinkerData(buf []byte) {
	//TODO 解析数据

	l := len(buf)
	crc := helper.ParseUint16(buf[l-2:])

	if crc != CRC16(buf[:l-2]) {
		//检验错误
		return
	}

	offset := helper.ParseUint16(buf[2:])
	length := 4
	switch buf[1] {
	case FuncCodeReadDiscreteInputs,
		FuncCodeReadCoils:
		count := int(helper.ParseUint16(buf[4:]))
		length += 1 + count/8
		if count%8 != 0 {
			length++
		}

		if l < length {
			//长度不够
			return
		}
		b := buf[6 : l-2]

		//解析开关
		bb := helper.ExpandBool(b, count)
		m.listener.OnAdapterRead(&models.Address{
			Slave:     buf[0],
			Offset:    offset,
			ReadCode:  buf[1],
		}, bb)

	case FuncCodeReadInputRegisters,
		FuncCodeReadHoldingRegisters,
		FuncCodeReadWriteMultipleRegisters:
		count := int(helper.ParseUint16(buf[4:]))
		length += 1 + count*2
		if l < length {
			//长度不够
			return
		}
		b := buf[6 : l-2]
		//解析word
		m.listener.OnAdapterRead(&models.Address{
			Slave:     buf[0],
			Offset:    offset,
			ReadCode:  buf[1],
		}, b)
	case FuncCodeWriteSingleCoil,
		FuncCodeWriteMultipleCoils,
		FuncCodeWriteSingleRegister,
		FuncCodeWriteMultipleRegisters:
		length += 4
		//TODO 处理写入成功
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

	//m.listener.OnAdapterRead()
	//m.listener.OnAdapterWrite()

}

func (m *RTU) OnLinkerError(err error) {
	m.listener.OnAdapterError(err)
}

func (m *RTU) OnLinkerClose() {}
