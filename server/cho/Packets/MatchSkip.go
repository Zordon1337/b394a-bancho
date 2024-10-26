package Packets

import (
	"net"
	"retsu/Utils"
)

func WriteMatchSkip(client net.Conn) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(62, []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
