package Packets

import (
	"net"
	"retsu/Utils"
)

func WriteMatchScoreUpdate(client net.Conn, data []byte) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(49, data) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
