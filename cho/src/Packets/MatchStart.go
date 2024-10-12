package Packets

import (
	"net"
	"socket-server/src/Structs"
	"socket-server/src/Utils"
)

func WriteMatchStart(client net.Conn, match Structs.Match) {
	resp, err := Utils.SerializePacket(47, Structs.GetBytesFromMatch(&match)) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
