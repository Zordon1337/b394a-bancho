package Packets

import (
	"net"
)

func WritePing(client net.Conn) {
	resp, err := SerializePacket(8, []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
