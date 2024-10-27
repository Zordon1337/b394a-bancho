package Packets

import (
	"net"
	"retsu/Utils"
)

func WriteMatchComplete(client net.Conn, build int) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(build, int(59))), []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
