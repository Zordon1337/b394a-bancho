package Packets

import (
	"net"
	"retsu/Utils"
	"retsu/cho/Structs"
)

func WriteMatchStart(client net.Conn, match Structs.Match) {
	if client == nil {
		return
	}
	resp, err := Utils.SerializePacket(47, Structs.GetBytesFromMatch(&match)) // empty packet lol
	if err != nil {
		return
	}
	client.Write(resp)
}
