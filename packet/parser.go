package packet

type Parser struct {
	buf []byte
}

func (p *Parser) Parse(buf []byte) []*Packet {
	//上次残存
	if p.buf != nil && len(p.buf) > 0 {
		buf = append(p.buf, buf...)
	}

	packs := make([]*Packet, 0)

	for {
		remain := len(buf)

		if remain < 6 {
			//包头都不够，等待剩余内容
			//可能需要 超时处理
			break
		}

		//寻找包头
		if buf[0] != '*' {
			buf = buf[1:]
			continue
		}
		if buf[1] != '#' {
			buf = buf[2:]
			continue
		}

		var msg Packet
		msg.Type = Type(buf[2])
		msg.Status = buf[3]

		l := int(uint16(buf[4])<<8 + uint16(buf[5]))
		if remain-6 < l {
			//内容不够，等待
			break
		}

		if l > 0 {
			msg.Data = dup(buf[6:])
		}

		packs = append(packs, &msg)

		//切片，继续解析
		buf = buf[6+l:]
	}

	if len(buf) > 0 {
		p.buf = dup(buf)
	}

	return packs
}

func dup(b []byte) []byte {
	buf := make([]byte, len(b))
	copy(buf, b)
	return buf
}
