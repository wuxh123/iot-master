package modbus

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Function Code
const (
	// Bit access
	FuncCodeReadDiscreteInputs = 2
	FuncCodeReadCoils          = 1
	FuncCodeWriteSingleCoil    = 5
	FuncCodeWriteMultipleCoils = 15

	// 16-bit access
	FuncCodeReadInputRegisters         = 4
	FuncCodeReadHoldingRegisters       = 3
	FuncCodeWriteSingleRegister        = 6
	FuncCodeWriteMultipleRegisters     = 16
	FuncCodeReadWriteMultipleRegisters = 23
	FuncCodeMaskWriteRegister          = 22
	FuncCodeReadFIFOQueue              = 24
	FuncCodeOtherReportSlaveID         = 17
	// FuncCodeDiagReadException          = 7
	// FuncCodeDiagDiagnostic             = 8
	// FuncCodeDiagGetComEventCnt         = 11
	// FuncCodeDiagGetComEventLog         = 12

)

// Exception Code
const (
	ExceptionCodeIllegalFunction                    = 1
	ExceptionCodeIllegalDataAddress                 = 2
	ExceptionCodeIllegalDataValue                   = 3
	ExceptionCodeServerDeviceFailure                = 4
	ExceptionCodeAcknowledge                        = 5
	ExceptionCodeServerDeviceBusy                   = 6
	ExceptionCodeNegativeAcknowledge                = 7
	ExceptionCodeMemoryParityError                  = 8
	ExceptionCodeGatewayPathUnavailable             = 10
	ExceptionCodeGatewayTargetDeviceFailedToRespond = 11
)


type address struct {
	unit uint8
	code uint8
	addr uint16
}

func (addr *address) toString () string  {
	return fmt.Sprintf("%d:%d:%d", addr.unit, addr.code, addr.addr)
}

// 设备号:功能码:地址
func parseAddress(addr string) (*address, error) {
	as := strings.Split(addr, ":")
	if len(as) < 3 {
		return nil, errors.New("地址格式必须是 设备号:功能码:地址")
	}

	var ad address
	v, e := strconv.Atoi(as[0])
	if e != nil {
		return nil, e
	}
	ad.unit = uint8(v)

	v, e = strconv.Atoi(as[1])
	if e != nil {
		return nil, e
	}
	ad.code = uint8(v)

	v, e = strconv.Atoi(as[2])
	if e != nil {
		return nil, e
	}
	ad.addr = uint16(v)

	return &ad, nil
}
