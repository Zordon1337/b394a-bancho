package Packets

import (
	"retsu/Utils"
	"retsu/cho/Structs"
)

func WriteMatchStart(client Structs.Player, match Structs.Match) {
	if client.Conn == nil {
		return
	}
	resp, err := Utils.SerializePacket(47, Structs.GetBytesFromMatch(&match, int32(client.Build))) // empty packet lol
	if err != nil {
		return
	}
	client.Conn.Write(resp)
}
