package modbus

import (
	"errors"
	"fmt"
	"git.zgwit.com/zgwit/iot-admin/base"
	"git.zgwit.com/zgwit/iot-admin/protocol/adapter"
	"git.zgwit.com/zgwit/iot-admin/protocol/helper"
)

func init() {
	adapter.RegisterAdapter(
		"Modbus RTU",
		adapter.Area{
			"线圈":    1,
			"触点":    2,
			"保持寄存器": 3,
			"输入寄存器": 4,
		},
		NewModbusRtu)
}

type RTU struct {
	link base.Link
}

func NewModbusRtu(opts string) (adapter.Adapter, error) {
	return &RTU{}, nil
}

func (m *RTU) Name() string {
	return "Modbus RTU"
}

func (m *RTU) Version() string {
	return "v0.0.1"
}

func (m *RTU) Read(slave uint8, area uint8, offset uint16, size uint16) ([]byte, error) {
	b := make([]byte, 8)
	b[0] = slave
	b[1] = area
	helper.WriteUint16(b[2:], offset)
	helper.WriteUint16(b[4:], size)
	helper.WriteUint16(b[6:], CRC16(b[:6]))

	buf, err := m.link.Request(b)
	if err != nil {
		return nil, err
	}

	//解析数据
	l := len(buf)
	crc := helper.ParseUint16(buf[l-2:])

	if crc != CRC16(buf[:l-2]) {
		//检验错误
		return nil, errors.New("校验错误")
	}

	//解析错误码
	if buf[1] & 0x80 > 0 {
		return nil, fmt.Errorf("错误码：%d", buf[2])
	}

	//解析数据
	length := 4
	count := int(helper.ParseUint16(buf[4:]))
	switch buf[1] {
	case FuncCodeReadDiscreteInputs,
		FuncCodeReadCoils:
		length += 1 + count/8
		if count%8 != 0 {
			length++
		}

		if l < length {
			//长度不够
			return nil, errors.New("长度不够")
		}
		b := buf[6 : l-2]
		//解析开关
		bb := helper.ExpandBool(b, count)
		return bb, nil
	case FuncCodeReadInputRegisters,
		FuncCodeReadHoldingRegisters,
		FuncCodeReadWriteMultipleRegisters:
		count := int(helper.ParseUint16(buf[4:]))
		length += 1 + count*2
		if l < length {
			//长度不够
			return nil, errors.New("长度不够")
		}
		b := buf[6 : l-2]
		return b, nil
	default:
		return nil, errors.New("不支持的指令")
	}
}

func (m *RTU) Write(slave uint8, area uint8, offset uint16, buf []byte) error {
	//如果是线圈，需要Shrink
	if area == 1 {
		buf = helper.ShrinkBool(buf)
		//TODO 长度需要计算
		//数据 转成 0x0000 0xFF00
	}

	l := 6 + len(buf)
	b := make([]byte, l)
	b[0] = slave
	b[1] = area //TODO 转写指令，并且考虑单个还是批量
	helper.WriteUint16(b[2:], offset)
	copy(b[4:], buf)
	helper.WriteUint16(b[l-2:], CRC16(b[:l-2]))

	buf, err := m.link.Request(b)
	if err != nil {
		return err
	}

	//解析数据
	l = len(buf)
	crc := helper.ParseUint16(buf[l-2:])

	if crc != CRC16(buf[:l-2]) {
		//检验错误
		return errors.New("校验错误")
	}

	//解析错误码
	if buf[1] & 0x80 > 0 {
		return fmt.Errorf("错误码：%d", buf[2])
	}

	return nil
}
