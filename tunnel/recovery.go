package tunnel

import (
	"iot-master/db"
	"iot-master/model"
	"log"
)

func Recovery() error {
	var ts []model.Tunnel
	err := db.Engine.Find(&ts)
	if err != nil {
		return err
	}

	for _, c := range ts {
		if c.Disabled {
			continue
		}
		_, err = StartTunnel(&c)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
