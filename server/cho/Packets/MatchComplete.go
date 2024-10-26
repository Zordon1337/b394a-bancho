package Packets

import (
	"net"
	"retsu/Utils"
)

func WriteMatchComplete(client net.Conn) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(59, []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
