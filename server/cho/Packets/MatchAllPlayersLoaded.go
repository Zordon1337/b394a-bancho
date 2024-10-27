package Packets

import (
	"net"
	"retsu/Utils"
)

func MatchAllPlayersLoaded(client net.Conn, build int) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(build, int(54))), []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
