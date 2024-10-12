package Packets

import (
	"net"
	"socket-server/src/Utils"
)

func WriteMatchComplete(client net.Conn) {
	resp, err := Utils.SerializePacket(59, []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
