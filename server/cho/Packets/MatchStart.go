package Packets

import (
	"retsu/Utils"
	"retsu/cho/Structs"
)

func WriteMatchStart(client Structs.Player, match Structs.Match, build int) {
	if client.Conn == nil {
		return
	}
	packet := 0
	if build > 1700 {
		packet = 48
	}
	if build < 1700 {
		packet = 48
	}
	resp, err := Utils.SerializePacket(int16(Utils.CalculatePacketOffset(int(build), int(packet))), Structs.GetBytesFromMatch(&match, int32(client.Build))) // empty packet lol
	if err != nil {
		return
	}
	client.Conn.Write(resp)
}
