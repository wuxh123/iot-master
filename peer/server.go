package peer

import (
	"github.com/zgwit/dtu-admin/packet"
	"log"
	"net"
)

func ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		peer := &Peer{
			conn:   conn,
			parser: packet.Parser{},
		}
		go peer.receive()
	}
}
