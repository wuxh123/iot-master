package modbus

import (
	"errors"
	"fmt"
	"iot-master/base"
	"iot-master/protocol"
	"iot-master/protocol/helper"
)

func init() {
	protocol.RegisterAdapter(
		"Modbus RTU",
		[]protocol.Code{
			{"线圈", 1},
			{"离散量", 2},
			{"保持寄存器", 3},
			{"输入寄存器", 4},
		},
		NewModbusRtu)
}

type response struct {
	buf []byte
	err error
}

type RTU struct {
	link base.Link
	resp chan response
}

func NewModbusRtu(opts string) (protocol.Adapter, error) {
	return &RTU{
		resp: make(chan response, 1),
	}, nil
}

func (m *RTU) Name() string {
	return "Modbus RTU"
}

func (m *RTU) Version() string {
	return "v0.0.1"
}

func (m *RTU) Attach(link base.Link) {
	m.link = link
	link.Listen(func(buf []byte) {

		//解析数据
		l := len(buf)
		crc := helper.ParseUint16(buf[l-2:])

		if crc != CRC16(buf[:l-2]) {
			//检验错误
			m.resp <- response{err: errors.New("校验错误")}
			return
		}

		//解析错误码
		if buf[1]&0x80 > 0 {
			m.resp <- response{err: fmt.Errorf("错误码：%d", buf[2])}
			return
		}

		//解析数据
		length := 4
		count := int(helper.ParseUint16(buf[1:]))
		switch buf[1] {
		case FuncCodeReadDiscreteInputs,
			FuncCodeReadCoils:
			length += 1 + count/8
			if count%8 != 0 {
				length++
			}

			if l < length {
				//长度不够
				m.resp <- response{err: errors.New("长度不够")}
				return
			}
			b := buf[2 : l-2]
			//数组解压
			bb := helper.ExpandBool(b, count)
			m.resp <- response{buf: bb}
		case FuncCodeReadInputRegisters,
			FuncCodeReadHoldingRegisters,
			FuncCodeReadWriteMultipleRegisters:
			length += 1 + count*2
			if l < length {
				//长度不够
				m.resp <- response{err: errors.New("长度不够")}
				return
			}
			b := buf[2 : l-2]
			m.resp <- response{buf: b}
		default:
			m.resp <- response{}
		}
	})
}

func (m *RTU) Read(slave uint8, code uint8, offset uint16, size uint16) ([]byte, error) {
	b := make([]byte, 8)
	b[0] = slave
	b[1] = code
	helper.WriteUint16(b[2:], offset)
	helper.WriteUint16(b[4:], size)
	helper.WriteUint16(b[6:], CRC16(b[:6]))

	err := m.link.Write(b)
	if err != nil {
		return nil, err
	}

	//等待结果
	resp := <-m.resp

	return resp.buf, resp.err
}

func (m *RTU) Write(slave uint8, code uint8, offset uint16, buf []byte) error {
	length := len(buf)
	//如果是线圈，需要Shrink
	if code == 1 {
		switch code {
		case FuncCodeReadCoils:
			if length == 1 {
				code = 5
				//数据 转成 0x0000 0xFF00
				if buf[1] > 0 {
					buf = []byte{0xFF, 0}
				} else {
					buf = []byte{0, 0}
				}
			} else {
				code = 15 //0x0F
				//数组压缩
				b := helper.ShrinkBool(buf)
				count := len(b)
				buf = make([]byte, 3+count)
				helper.WriteUint16(buf, uint16(length))
				buf[2] = uint8(count)
				copy(buf[3:], b)
			}
		case FuncCodeReadHoldingRegisters:
			if length == 2 {
				code = 6
			} else {
				code = 16 //0x10
				b := make([]byte, 3+length)
				helper.WriteUint16(b, uint16(length/2))
				b[2] = uint8(length)
				copy(b[3:], buf)
				buf = b
			}
		}
	}

	l := 6 + len(buf)
	b := make([]byte, l)
	b[0] = slave
	b[1] = code
	helper.WriteUint16(b[2:], offset)
	copy(b[4:], buf)
	helper.WriteUint16(b[l-2:], CRC16(b[:l-2]))

	err := m.link.Write(b)
	if err != nil {
		return err
	}

	//等待结果
	resp := <-m.resp

	return resp.err
}
