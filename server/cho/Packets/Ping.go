package Packets

import (
	"net"
	"retsu/Utils"
)

func WritePing(client net.Conn) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(8, []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
