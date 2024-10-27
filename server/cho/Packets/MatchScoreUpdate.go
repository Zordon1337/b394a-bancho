package Packets

import (
	"net"
	"retsu/Utils"
)

func WriteMatchScoreUpdate(client net.Conn, data []byte, build int) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(int(build), int(49))), data) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
