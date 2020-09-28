package modbus

func LRC(buf []byte) uint8 {
	var sum uint8 = 0
	for _, b := range buf {
		sum += b
	}
	return uint8(-int8(sum))
}
