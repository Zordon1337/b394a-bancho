package Packets

import (
	"net"
	"retsu/Utils"
)

func WriteMatchSkip(client net.Conn, build int) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(int(build), int(62))), []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
