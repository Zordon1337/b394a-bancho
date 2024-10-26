package Packets

import (
	"retsu/Utils"
	"retsu/cho/Structs"
)

func WriteMatchUpdate(client Structs.Player, match Structs.Match) {
	if client.Conn == nil {
		return
	}
	bytes := Structs.GetBytesFromMatch(&match, int32(client.Build))
	resp, err := Utils.SerializePacket(27, bytes)
	if err != nil {
		return
	}
	client.Conn.Write(resp)
}
