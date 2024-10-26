package Packets

import (
	"net"
	"retsu/Utils"
	"retsu/cho/Structs"
)

func WriteMatchJoinSuccess(client net.Conn, match Structs.Match, build int32) {
	if client == nil {
		return
	}

	resp, err := Utils.SerializePacket(37, Structs.GetBytesFromMatch(&match, build))
	if err != nil {
		return
	}
	client.Write(resp)
}
