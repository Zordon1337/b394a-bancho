package Packets

import (
	"net"
	"socket-server/src/Utils"
)

func WritePing(client net.Conn) {
	resp, err := Utils.SerializePacket(8, []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
