package Packets

import (
	"net"
	"socket-server/src/Utils"
)

func MatchAllPlayersLoaded(client net.Conn) {
	resp, err := Utils.SerializePacket(54, []byte{}) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
