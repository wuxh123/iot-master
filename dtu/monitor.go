package dtu

import (
	"encoding/hex"
	"github.com/gorilla/websocket"
	"time"
)

type monitorPack struct {
	Type string    `json:"type"`
	Time time.Time `json:"time"`
	Data string    `json:"data"`
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
		Time: time.Now(),
		Data: da,
	})
}

func (m *Monitor) Receive() {
	defer m.Conn.Close()
	for {
		var msg monitorPack
		//读取ws中的数据
		err := m.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		switch msg.Type {
		case "send":
			b, e := hex.DecodeString(msg.Data)
			if e != nil {
				break
			}
			m.Link.Send(b)
		case "ping":
			m.Conn.WriteJSON(&monitorPack{
				Type: "pong",
			})
		}
	}

	m.Link.monitor = nil
}
