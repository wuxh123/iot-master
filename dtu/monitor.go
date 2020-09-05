package dtu

import (
	"encoding/hex"
	"github.com/gorilla/websocket"
)

type monitorPack struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Monitor struct {
	Key  string
	Conn *websocket.Conn
	Link *Link
}

func (m *Monitor) Report(typ string, data []byte) error {
	da := ""
	if data != nil {
		da = hex.EncodeToString(data)
	}
	return m.Conn.WriteJSON(&monitorPack{
		Type: typ,
		Data: da,
	})
}

func (m *Monitor) receive() {
	defer m.Conn.Close()
	for {
		var msg monitorPack
		//读取ws中的数据
		err := m.Conn.ReadJSON(&msg) //.ReadMessage()
		if err != nil {
			break
		}

		switch msg.Type {
		//TODO 由接口来发
		case "transfer":
			b, _ := hex.DecodeString(msg.Data)
			if b != nil {
				_, _ = m.Link.Send(b)
			}
		}
	}

	m.Link.monitors.Delete(m)
}
