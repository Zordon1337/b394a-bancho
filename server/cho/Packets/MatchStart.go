package Packets

import (
	"retsu/Utils"
	"retsu/cho/Structs"
)

func WriteMatchStart(client Structs.Player, match Structs.Match, build int) {
	if client.Conn == nil {
		return
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(int(build), int(47))), Structs.GetBytesFromMatch(&match, int32(client.Build))) // empty packet lol
	if err != nil {
		return
	}
	client.Conn.Write(resp)
}
