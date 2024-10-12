package Packets

import (
	"net"
	"socket-server/src/Utils"
)

func WriteMatchScoreUpdate(client net.Conn, data []byte) {
	resp, err := Utils.SerializePacket(49, data) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
