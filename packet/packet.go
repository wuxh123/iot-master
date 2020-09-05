package packet


type Packet struct {
	Type   Type
	Status uint8
	Data   []byte
}

func (p *Packet) Encode() []byte  {
	l := 6
	ln := 0
	if p.Data != nil {
		ln = len(p.Data)
	}
	buf := make([]byte, l + ln)
	buf[0] = '*'
	buf[1] = '#'
	buf[2] = uint8(p.Type)
	buf[3] = p.Status
	buf[4] = uint8(ln >> 8)
	buf[5] = uint8(ln)
	if ln > 0 {
		copy(buf[6:], p.Data)
	}
	return buf
}

